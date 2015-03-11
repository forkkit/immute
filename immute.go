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
	Obj() interface{}
	Add(f interface{}, i interface{})
	Get(f interface{}) interface{}
	Delete(f interface{}) interface{}
	Clear()
}

//ListSequence for all list related collections,its the root of all list sequences
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

func (s *ListSequence) Clear() {
	s.Data = make([]interface{}, 0)
}

func (s *ListSequence) Push(i interface{}) {
	s.Add(i, nil)
}

func (s *ListSequence) Add(i interface{}, _ interface{}) {
	s.Data = append(s.Data, i)
}

func (s *ListSequence) Delete(f interface{}) interface{} {
	index, ok := f.(int)

	if !ok {
		return nil
	}

	// rd := *s.Data
	// rd = append(rd[:index], rd[:index+1]...)
	item := s.Data[index]
	s.Data = append(s.Data[:index], s.Data[index+1:]...)
	// s.Data = *rd
	return item
}

func (s *ListSequence) Get(f interface{}) interface{} {
	ind, ok := f.(int)
	if !ok {
		return nil
	}
	return s.Data[ind]
}

func (s *ListSequence) Set(f interface{}, i interface{}) {
	ind, ok := f.(int)

	if !ok {
		return
	}

	if ind >= len(s.Data) || ind < 0 {
		return
	}

	s.Data[ind] = i
}

func (s *ListSequence) Length() int {
	return len(s.Data)
}

func (s *ListSequence) Obj() interface{} {
	return s.Data
}

func (s *ListSequence) Seq() *Sequence {
	return &Sequence{Sequencer(s)}
}

//MapSequence for all map related collections,its the root of all map sequences
type MapSequence struct {
	Data map[interface{}]interface{}
}

func (s *MapSequence) Clear() {
	s.Data = make(map[interface{}]interface{})
}

func (s *MapSequence) Get(f interface{}) interface{} {
	return s.Data[f]
}

func (s *MapSequence) Set(f interface{}, i interface{}) {
	s.Data[f] = i
}

func (s *MapSequence) Add(i interface{}, f interface{}) {
	s.Set(i, f)
}

func (s *MapSequence) Delete(i interface{}) interface{} {
	item := s.Data[i]
	delete(s.Data, i)
	return item
}

func (s *MapSequence) Obj() interface{} {
	return s.Data
}

func (s *MapSequence) Length() int {
	return len(s.Data)
}

func (s *MapSequence) Seq() *Sequence {
	return &Sequence{Sequencer(s)}
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

//Sequence is the core atomic structure for all sequence based operations
type Sequence struct {
	Parent Sequencer
}

func (s *Sequence) Delete(i interface{}) interface{} {
	return s.Parent.Delete(i)
}

func (s *Sequence) Clear() {
	s.Parent.Clear()
}

func (s *Sequence) Get(i interface{}) {
	s.Parent.Get(i)
}

func (s *Sequence) Add(i interface{}, f interface{}) {
	s.Parent.Add(i, f)
}

func (s *Sequence) Each(f EachHandler, c CompleteHandler) *Sequence {
	return s.Parent.Each(f, c)
}

func (s *Sequence) Obj() interface{} {
	return s.Parent.Obj()
}

func (s *Sequence) Length() int {
	return s.Parent.Length()
}

//SequenceOp is the core of all sequence operations,you define a sequence op that
//returns a new sequence as its result
type SequenceOp struct {
	Root       *Sequence
	ParentEach func(r *Sequence, f EachHandler, c CompleteHandler) *Sequence
	EachItem   EachHandler
	Completed  CompleteHandler
}

//MemoizedSequenceOp is a sequence operation to memoize your sequence operation for fast retrieval without re-doing all the work
type MemoizedSequenceOp struct {
	op    *SequenceOp
	Cache *Sequence
}

func (s *MemoizedSequenceOp) Each() *Sequence {
	if s.Cache != nil {
		return s.Cache
	}
	s.Cache = s.op.Each().Seq()
	return s.Cache
}

func (s *MemoizedSequenceOp) Length() interface{} {
	return s.Each().Length()
}

//Memoize allows caching of the operation of the current sequenceOp
func (s *SequenceOp) Memoize() interface{} {
	return &MemoizedSequenceOp{s, nil}
}

//Length returns the total size of the collection
func (s *SequenceOp) Length() interface{} {
	return s.Each().Length()
}

//Each iterates through all collections
func (s *SequenceOp) Each() *Sequence {
	return s.ParentEach(s.Root, s.EachItem, s.Completed)
}

//Map provides a mutating of sequence values by a function
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

//Filter provides a means of filtering data within the collection into a new sequence
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

func (s *MemoizedSequenceOp) Map(fe EachHandler, co CompleteHandler) *SequenceOp {
	return s.Each().Map(fe, co)
}

func (s *SequenceOp) Map(fe EachHandler, co CompleteHandler) *SequenceOp {
	return s.Each().Map(fe, co)
}

func (s *Sequence) Map(fe EachHandler, co CompleteHandler) *SequenceOp {
	return &SequenceOp{s, Map, fe, co}
}

func (s *MemoizedSequenceOp) Filter(fe EachHandler, co CompleteHandler) *SequenceOp {
	return s.Each().Filter(fe, co)
}

func (s *SequenceOp) Filter(fe EachHandler, co CompleteHandler) *SequenceOp {
	return s.Each().Filter(fe, co)
}

func (s *Sequence) Filter(fe EachHandler, co CompleteHandler) *SequenceOp {
	return &SequenceOp{s, Filter, fe, co}
}

func (s *SequenceOp) Obj() interface{} {
	return s.Each().Obj()
}

func (s *SequenceOp) Seq() *Sequence {
	return s.Each()
}

func (s *Sequence) Seq() *Sequence {
	return s
}

func CreateList(i []interface{}) *Sequence {
	return (&ListSequence{i}).Seq()
}

func CreateMap(i map[interface{}]interface{}) *Sequence {
	return (&MapSequence{i}).Seq()
}

//CreateListSeq creates a pure list sequence  of any type
func CreateListSeq(i []interface{}) *ListSequence {
	return &ListSequence{i}
}

//CreateMapSeq creates a pure map sequence  of any type
func CreateMapSeq(i map[interface{}]interface{}) *MapSequence {
	return &MapSequence{i}
}
