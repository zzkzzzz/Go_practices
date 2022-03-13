# Notes

## Some Tips

- `defer` execution order is LIFO
  - The defer statement adds the function call following the defer keyword onto a stack.
  Since the calls are placed on a stack, they are called in last-in-first-out (LIFO) order.

  ```go
  for i := 0; i < 5; i++ {
    defer fmt.Printf("%d", i)
  }
  ```

- Passing by value
  - In Go, there is no such thing as passing by reference. Everything is passed by value. If you assign the value of an array to another variable, the entire value is copied.

- Slice Memory Representation
  - A slice, unlike an array, does not allocate the memory of the data blocks during initialization. In fact, slices are initialized with the nil value.
  - A slice is allocated differently from an array, and is actually a modified pointer. Each slice contains three pieces of information:
  
  ```go
  type slice struct {
    array unsafe.Pointer
	  len   int
	  cap   int
  }
  ```
  - reference passing
    - If you simply assign an existing slice to a new variable, the slice will not be duplicated.
      - When you assign a slice to another variable, you still pass by value. The value here refers to just the pointer, length, and capacity, and not the memory occupied by the elements themselves.
      - Both variables will refer to the exact same slice, so any changes to the slice’s value will be reflected in all its references. This is also true when passing a slice to a function since slices are passed by reference in Go.
      - try to use builtin `copy` or `append` functions
      - perfermance of `copy` is better than `append` when copy slice
  - slice `append` is not thread-safe
  