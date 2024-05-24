package myutils

import (
	"cmp"
	"fmt"
	"slices"
	"sync"
	"time"
)

const EVENT_CHAIN_SIZE int = 25
const MIN_EVENT_TIME = 1 // Any event with less than 1μs time elapsed or blank name is not reported.

type TimerEvent struct {
	name  string
	clock int64
}

func NewTimerEvent(name string) *TimerEvent {
	var now = time.Now().UTC()
	event := &TimerEvent{
		name:  name,
		clock: now.UnixMicro(),
	}
	return event
}

/*
* The Event History "class" to aggregate the events logged to the GroupTimer.
* Aggregation occurs when we have filled the EventChain
 */
type EventHistory struct {
	name  string
	start int64  // EPoch time in micro seconds
	total int64  // Total time in micro seconds
	hits  uint32 // Number of time this event occured
}

func NewEventHistory(name string, clock int64) *EventHistory {
	result := &EventHistory{
		name:  name,
		start: clock,
		total: 0,
		hits:  0,
	}
	return result
}

func (eh *EventHistory) eventHistoryToString() string {
	if eh.hits != 1 {
		return fmt.Sprintf("%v:%v hits=%v", eh.name, eh.total, eh.hits)
	}
	return fmt.Sprintf("%v:%v", eh.name, eh.total)
}

func (h *EventHistory) addEventDuration(duration int64) {
	h.hits++
	h.total += duration
}

/*
* The GroupTimer "class" supporting a collection of timed events (event has a name and start time)
* By logging successive events we can calculate durations.
*
* GroupTimer aggregates the time duration for each named events
* All events finsih time is based on the subsequent event.
* A blank name event is used to act as "just an end marker" for the previous event.
 */
type GroupTimer struct {
	mutex      *sync.Mutex
	group      int
	firstClock int64
	lastClock  int64
	eventChain []TimerEvent
	eventCount int
	historyMap map[string]*EventHistory
}

func NewGroupTimer(grp int) *GroupTimer {
	var gt = &GroupTimer{
		mutex:      &sync.Mutex{},
		group:      grp,
		firstClock: 0,
		lastClock:  0,
		eventChain: make([]TimerEvent, EVENT_CHAIN_SIZE),
		eventCount: 0,
		historyMap: make(map[string]*EventHistory),
	}
	//	for i := 0; i < EVENT_CHAIN_SIZE; i++ {
	//		gt.eventChain[i] = TimerEvent{} // seed with empty timer events (is this really needed?)
	//	}
	return gt
}

func (gt *GroupTimer) log(name string) *GroupTimer {
	var event = NewTimerEvent(name)
	if gt.firstClock == 0 {
		gt.firstClock = event.clock
	}
	gt.lastClock = event.clock
	gt.mutex.Lock()
	defer gt.mutex.Unlock()
	// Now add the event to the eventChain
	gt.eventChain[gt.eventCount] = *event
	gt.eventCount++
	if gt.eventCount >= EVENT_CHAIN_SIZE {
		gt.compact()
	}
	return gt
}

// This must never be called without the gt.mutex lock in place!!
func (gt *GroupTimer) compact() *GroupTimer {
	// TODO add in the compact code moving the eventChain to the aggregated EventHistory(s)
	if gt.eventCount == 0 {
		return gt
	}
	currentEvent := gt.eventChain[0]
	for i := 1; i < gt.eventCount; i++ {
		nextEvent := gt.eventChain[i]
		duration := nextEvent.clock - currentEvent.clock
		history, found := gt.historyMap[currentEvent.name]
		if !found {
			history = NewEventHistory(currentEvent.name, currentEvent.clock)
			gt.historyMap[currentEvent.name] = history
		}
		history.addEventDuration(duration)
		currentEvent = nextEvent
	}
	// Reset the event chain to contain the last event only
	// Note until we get another event we do not know the duration of this currentEvent
	// That is why the phaseTimer has done() and pause() methods (each equivalent) to provide end clock times
	gt.eventCount = 0
	gt.eventChain[0] = currentEvent
	return gt
}

func (gt *GroupTimer) groupTimerToString() string {
	gt.mutex.Lock()
	defer gt.mutex.Unlock()
	var res = fmt.Sprintf("G%v %vμs", gt.group, gt.lastClock-gt.firstClock)
	gt.compact()
	for eventName, value := range gt.historyMap {
		if eventName == "" || value.total < MIN_EVENT_TIME {
			continue
		}
		res = fmt.Sprintf("%v {%v}", res, value.eventHistoryToString())
	}
	return res
}

/*
**** PhaseTimer - the timer uses groups to separate out timing of different event streams. For example, during
any thread processing where threads can share the same timer. Unlike Java - Go has no thread Id that can be identified.
This puts the onus on the Go caller to use groups as a way to distinquish between different threads.
Of course using separate timers would also work... but is less clear when it comes to looking at the timer outputs.
*/
type PhaseTimer struct {
	phaseMutex  *sync.Mutex
	groupTimers map[int]*GroupTimer
}

func NewPhaseTimer() *PhaseTimer {
	var pt = &PhaseTimer{
		phaseMutex:  &sync.Mutex{},
		groupTimers: make(map[int]*GroupTimer),
	}
	return pt
}

/* Internal function
 */
func (pt *PhaseTimer) getOrAddGroupTimer(groupId int) *GroupTimer {
	pt.phaseMutex.Lock()
	defer pt.phaseMutex.Unlock()
	groupTimer, found := pt.groupTimers[groupId]
	if !found {
		groupTimer = NewGroupTimer(groupId)
		pt.groupTimers[groupId] = groupTimer
	}
	return groupTimer
}

/** Simple log all items are logged into group 0
 */
func (pt *PhaseTimer) LogN(groupId int, name string) *PhaseTimer {
	groupTimer := pt.getOrAddGroupTimer(groupId)
	groupTimer.log(name)
	return pt
}

func (pt *PhaseTimer) DoneN(groupId int) *PhaseTimer {
	return pt.LogN(groupId, "")
}

// Convenience functions for simple (non group N) logging of events
func (pt *PhaseTimer) Log(name string) *PhaseTimer {
	return pt.LogN(0, name)
}

func (pt *PhaseTimer) Done() *PhaseTimer {
	return pt.DoneN(0)
}

func (pt *PhaseTimer) AllDone() *PhaseTimer {
	for _, timer := range pt.groupTimers {
		timer.log("")
	}
	return pt
}

// Comparator for GroupTimers:
func cmpTimer(a, b *GroupTimer) int {
	return cmp.Compare(a.firstClock, b.firstClock)
}

func (pt *PhaseTimer) getSortedGroupTimers() []*GroupTimer {
	timerArr := make([]*GroupTimer, len(pt.groupTimers))
	count := 0
	for _, timer := range pt.groupTimers {
		timerArr[count] = timer
		count++
	}
	slices.SortFunc(timerArr, cmpTimer)
	return timerArr
}

func (pt *PhaseTimer) ToString() string {
	res := "Timings:"
	if len(pt.groupTimers) == 0 {
		return res + "None"
	}
	timerArr := pt.getSortedGroupTimers()
	for _, timer := range timerArr {
		res = fmt.Sprintf("%v [%v]", res, timer.groupTimerToString())
	}

	return res
}
