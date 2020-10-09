// automatically generated by stateify.

package futex

import (
	"gvisor.dev/gvisor/pkg/state"
)

func (b *bucket) StateTypeName() string {
	return "pkg/sentry/kernel/futex.bucket"
}

func (b *bucket) StateFields() []string {
	return []string{}
}

func (b *bucket) beforeSave() {}

func (b *bucket) StateSave(stateSinkObject state.Sink) {
	b.beforeSave()
	if !state.IsZeroValue(&b.waiters) {
		state.Failf("waiters is %#v, expected zero", &b.waiters)
	}
}

func (b *bucket) afterLoad() {}

func (b *bucket) StateLoad(stateSourceObject state.Source) {
}

func (m *Manager) StateTypeName() string {
	return "pkg/sentry/kernel/futex.Manager"
}

func (m *Manager) StateFields() []string {
	return []string{
		"sharedBucket",
	}
}

func (m *Manager) beforeSave() {}

func (m *Manager) StateSave(stateSinkObject state.Sink) {
	m.beforeSave()
	if !state.IsZeroValue(&m.privateBuckets) {
		state.Failf("privateBuckets is %#v, expected zero", &m.privateBuckets)
	}
	stateSinkObject.Save(0, &m.sharedBucket)
}

func (m *Manager) afterLoad() {}

func (m *Manager) StateLoad(stateSourceObject state.Source) {
	stateSourceObject.Load(0, &m.sharedBucket)
}

func (l *waiterList) StateTypeName() string {
	return "pkg/sentry/kernel/futex.waiterList"
}

func (l *waiterList) StateFields() []string {
	return []string{
		"head",
		"tail",
	}
}

func (l *waiterList) beforeSave() {}

func (l *waiterList) StateSave(stateSinkObject state.Sink) {
	l.beforeSave()
	stateSinkObject.Save(0, &l.head)
	stateSinkObject.Save(1, &l.tail)
}

func (l *waiterList) afterLoad() {}

func (l *waiterList) StateLoad(stateSourceObject state.Source) {
	stateSourceObject.Load(0, &l.head)
	stateSourceObject.Load(1, &l.tail)
}

func (e *waiterEntry) StateTypeName() string {
	return "pkg/sentry/kernel/futex.waiterEntry"
}

func (e *waiterEntry) StateFields() []string {
	return []string{
		"next",
		"prev",
	}
}

func (e *waiterEntry) beforeSave() {}

func (e *waiterEntry) StateSave(stateSinkObject state.Sink) {
	e.beforeSave()
	stateSinkObject.Save(0, &e.next)
	stateSinkObject.Save(1, &e.prev)
}

func (e *waiterEntry) afterLoad() {}

func (e *waiterEntry) StateLoad(stateSourceObject state.Source) {
	stateSourceObject.Load(0, &e.next)
	stateSourceObject.Load(1, &e.prev)
}

func init() {
	state.Register((*bucket)(nil))
	state.Register((*Manager)(nil))
	state.Register((*waiterList)(nil))
	state.Register((*waiterEntry)(nil))
}
