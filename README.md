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
    
    accessing module function or class is like accessing class member
    example:
      module_name.function
      var_name:module_name.class
      var_name = module_name.function
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
    
  # Variables
    Variables may be declare
      a:string
    Or just assigned
      a = 10
  #### They are type safe, once declared or assigned, sets the type
    
  # Functions
    name() {
    
    }
    
    name(n:string,n:number) {
      println("ok")
    }
    
    name(n:bool):bool {
      println('done')
      return false
    }
  
  # Classes definition
  
  ### there is no keyword for class
  The compiler recognizes a class construction by :
  
  name {
  
  }
  
  #### main is the exception
  
  ## Contructor
  ### this
    this(n:number)
    may have arguments like functions, but will not permit simple declaration except on function return value 
    may declare functions with the same name, but diferent parameters count or type, like java or c++
    
 ### Fields, Methods
 If starts by a single _ means its protected. Still work to be done here
 If starts with two _ means private
 
