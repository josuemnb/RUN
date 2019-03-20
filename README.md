# Run
  Nice programming language inspired in golang with classes and generics like java
  Compiled and fast as C
  
  Example:
  
    import files
    
    class {
      f():string {
        return "Hello World"
      }
    }

    func():number {
      return (10*13/4)
    }

    main {
      info = files.info('test.run')
      println(info.exist())
      println(func())
      c:class
      println(c.f())
    }

# Documentation
  ## Basic data types
  bool, number, real, string (string works as class)
  
  ## Builtin Functions
  ### print(ln)
    print(ln)(string|bool|number|real)
    behaves like golang println - values separeted by comma, but without space between them
  
  ## Keywords
  
  ### import
    import modules
  ### module
    first word inside modules
  ### main
    keyword that defines main program - needs to be present in order to compile
    it can be defined without paranteses but can accepts parameters as a function
  ### loop
    similar to for/while 
    example : 
              
              loop a=0..1000 { //if variable doesnt exist in scopes, its created
                println(a)
              }
              
              loop ..1000 { // loop value from 0 to 1000
              
              }
              
              loop { // forever
              
              }
              
              loop (bool) { //behaves like while
              
              }
              
  ### return, break, nil, false, true
    like any other language
  ### this
    refering to the current class construction
  ### if else 
    conditional like golang - without parenteses
  ### cpp 
    special keyword to inject c code
    very helpful while creating modules
  
  
