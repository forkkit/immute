#Immute.go
 Provides lazy sequences for golang and a uniform iteraterable and modifable opertions on maps or list of items

##Install

      go get github.com/influx6/immute

  Then

      go install github.com/influx6/immute

##Example

  ```

      import (
        "github.com/influx6/immute"
      )

      seq := immute.CreateMap(map[interface{}]interface{}{1: "sic", 3: "luc"})


      filter := seq.Filter(func(i interface{}, k interface{}) interface{} {
        return i.(string) == "sic"
      }, func(c int, f interface{}) {})

      map := seq.Map(func(i interface{}, k interface{}) interface{} {
        return (k.(string) + i.(int))
      }, func(c int, f interface{}) {
          // notify me when we done iterating
      })

      //get the sequence result lazily and it generates the sequence when you call this
      filtered := filter.Obj()
      mapped := map.Obj()


  ```

##API
    
####Sequences
  Early on immute was about sequences,a set of mutable values that can be iterated ,map and performed operation on, although at first glance golang may seem not to suite such dynamic behaviour but it is all in the eyes that sees it,ofcourse a few things had to be adopted to allow the global applyable sequence operations without having to go to reflection for that.

  -  Sequencer
      Both List and Map sequences or any other sequence must match the `Sequener` interface which has the follow methods:

  - Sequencer.Each(each func(value, key interface{}), completed func(length int, object interface{}))
      This method allows the iteration over all values of a sequence,its the core behind all `Sequence` based operations such as filter,map,take,...etc. All sequencable must implement this.

  - Sequencer.Length() int
      This method returns the length of the sequencable

  - Sequencer.Add(key,value interface{})
      This method is for adding a new item into the sequence but in the case of the ListSequence the second should be set to nil as its not needed nor used

  - Sequencer.Get(key interface{}) interface{}
     This method used to get a value by a key from the sequence

  - Sequencer.Delete(key interface{}) 
     This method remove an item by its key from the sequence

  - Sequencer.Clear()
     This method basically clears and empties the sequence
                
  - ListSequence
     This is the definition for list/array based sequence types. The structs definies that its data be a  []interface{} type to allow any data type to be added and immute a function for encapsulating all the details for easy use.

            data := immute.CreateList([]interface{1,4,6,7,8,9,10})
        
            data.Seq().Each(...)
            data.Seq().Map(...).filter(...)

  - MapSequence
    This is the definition for map based sequence types. The structs definies that its data be a  map[interface{}]interface{} type to allow any data type to be added and immute a function for encapsulating all the details for easy use.

            data := immute.CreateList(map[interface{}]interface{‘name”:”ally”,”tel”:07087723232 })
        
            data.Seq().Each(...)
            data.Seq().Map(...).filter(...)



####SequenceOperation (SequenceOp)
 These are sequencable processors,they take sequences to create another sequence from their original versions and allows the basic operations, among such include:

   - Filter: this create a  `SequenceOp` that filters out the sequence based on a predicate provided 
      
      ```
            data := immute.CreateList(map[interface{}]interface{‘name”:”ally”,”tel”:07087723232 })
        
            filtered := data.Filter(func (val interface{},key interface{}) interface{} {
                return true/false
            },func (total int, obj interface{}){
                //its done
            })
     ```

- Map: this create a  `SequenceOp` that maps out the result sequence on a predicate provided 
      
      ```
            data := immute.CreateList(map[interface{}]interface{‘name”:”ally”,”tel”:07087723232 })
        
            mapped := data.Map(func (val interface{},key interface{}) interface{} {
                return ....
            },func (total int, obj interface{}){
                //its done
            })
     ```


####Immutable Structures 
 These is the main goal for immute as a package although this has not being implemented yet but its the creation of deeply nestable immutable structures like those in the Facbook reactjs framework which allow more than must notification but time reversal of operations  (to be implemented....)



