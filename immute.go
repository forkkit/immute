/*
  Immute provides sequences and immutable structures for go, its the cross-over
  implementation of the immutable library from the js library Stackq
  Stackq (https://github.com/influx6/stackq)

*/

package immute

type EachHandler func(v interface{}, f interface{}) interface{}
type CompleteHandler func(r int, a interface{})

type Sequencer interface {
	Each(f EachHandler, c CompleteHandler) *Sequence
	Length() int
	toObj() interface{}
}

type ListSequence struct {
	Data []interface{}
}

func (s *ListSequence) Each(f EachHandler, c CompleteHandler) *Sequence {
	count := 0

	for key, value := range s.Data {
		f(value, key)
		count += 1
	}
	c(count, s)
	return nil
}

func (s *ListSequence) Length() int {
	return len(s.Data)
}

func (s *ListSequence) toObj() interface{} {
	return s.Data
}

type MapSequence struct {
	Data map[interface{}]interface{}
}

func (s *MapSequence) toObj() interface{} {
	return s.Data
}

func (s *MapSequence) Length() int {
	return len(s.Data)
}

func (s *MapSequence) Each(f EachHandler, c CompleteHandler) *Sequence {
	count := 0

	for key, value := range s.Data {
		f(value, key)
		count += 1
	}
	c(count, s)
	return nil
}

func CreateListSeq(i []interface{}) *ListSequence {
	return &ListSequence{i}
}

func CreateMapSeq(i map[interface{}]interface{}) *MapSequence {
	return &MapSequence{i}
}

type Sequence struct {
	Parent Sequencer
	Size   int
}

func (s *Sequence) Each(f EachHandler, c CompleteHandler) *Sequence {
	return s.Parent.Each(f, c)
}

func (s *Sequence) toObj() interface{} {
	return s.Parent.toObj()
}

func (s *Sequence) Length() interface{} {
	return s.Parent.Length()
}

type SequenceOp struct {
	Root       *Sequence
	ParentEach func(r *Sequence, f EachHandler, c CompleteHandler) *Sequence
	EachItem   EachHandler
	Completed  CompleteHandler
}

type MemoizedSequenceOp struct {
	op    *SequenceOp
	Cache *Sequence
}

func (s *MemoizedSequenceOp) Each() *Sequence {
	if s.Cache != nil {
		return s.Cache
	}
	s.Cache = s.op.Each()
	return s.Cache
}

func (s *MemoizedSequenceOp) Length() int {
	return this.Each().Length()
}

func (s *SequenceOp) Memoize() interface{} {
	return &MemoizedSequenceOp{s}
}

func (s *SequenceOp) Length() int {
	return s.Each().Length()
}

func (s *SequenceOp) Each() *Sequence {
	return s.ParentEach(s.Root, s.EachItem, s.Completed)
}

func Map(s *Sequence, feach EachHandler, comp CompleteHandler) *Sequence {
	count := s.Parent.Length()
	data := make([]interface{}, 0)
	_ = s.Parent.Each(func(i interface{}, r interface{}) interface{} {
		conv := feach(i, r)

		if conv != nil {
			data = append(data, conv)
		}

		return count
	}, func(c int, _ interface{}) {
		comp(c, data)
	})

	return CreateList(data)
}

func Filter(s *Sequence, feach EachHandler, comp CompleteHandler) *Sequence {
	count := 0
	data := make([]interface{}, 0)
	_ = s.Parent.Each(func(i interface{}, r interface{}) interface{} {
		conv := feach(i, r)
		state, ok := conv.(bool)

		if !ok {
			return count
		}

		if state {
			data = append(data, i)
			count += 1
		}

		return count
	}, func(c int, _ interface{}) {
		comp(c, data)
	})

	return CreateList(data)
}

func (s *SequenceOp) Map(fe EachHandler, co CompleteHandler) *SequenceOp {
	return s.Each().Map(fe, co)
}

func (s *Sequence) Map(fe EachHandler, co CompleteHandler) *SequenceOp {
	return &SequenceOp{s, Map, fe, co}
}

func (s *SequenceOp) Filter(fe EachHandler, co CompleteHandler) *SequenceOp {
	return s.Each().Filter(fe, co)
}

func (s *Sequence) Filter(fe EachHandler, co CompleteHandler) *SequenceOp {
	return &SequenceOp{s, Filter, fe, co}
}

func (s *SequenceOp) toObj() interface{} {
	return s.Each().toObj()
}

func CreateList(i []interface{}) *Sequence {
	l := &ListSequence{i}
	return &Sequence{Sequencer(l), len(i)}
}

func CreateMap(i map[interface{}]interface{}) *Sequence {
	m := &MapSequence{i}
	return &Sequence{Sequencer(m), len(i)}
}
