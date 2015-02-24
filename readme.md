#Immute.go
 Provides lazy sequences for golang and a uniform iteraterable and modifable opertions on maps or list of items

##Install
  go get http://github.com/influx6/immute


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
