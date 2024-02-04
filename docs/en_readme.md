## Introduction
Kahroba (amber stone) is a simple and flexible language for learning purposes.
It'll help you learn how a interepreted language works.

## How to use
First of all, make a file with .krb or .kahroba extention, then run interpreter from command line:
```bash
./kahroba main.kahroba     // linux
kahroba.exe main.kahroba   // windows
kahroba_mac main.kahroba   // mac
``` 

Install Suitable version from links below:

[Linux](https://github.com/kahroba-lang/kahroba/releases/download/0.1/kahroba) \
[Windows](https://github.com/kahroba-lang/kahroba/releases/download/0.1/kahroba.exe) \
[MacOS](https://github.com/kahroba-lang/kahroba/releases/download/0.1/kahroba_mac)

Also, you can build it from source, but you'll need go v1.19:
```bash
$ git clone https://github.com/kahroba-lang/kahroba.git
$ cd kahroba
$ go build
```

## Comments

Like most of languages, you can use // to define your comments, they won't get interpreted:
```rust
// This is my first program in Kahroba programming language, Let's Rock!
```

## Strings
You can define string using double quotation:
```rust
"Hello World!"
```
They can concatenate using + operator:
```rust
"Hello " + "World!" // Hello World
```

Right side type of an operations automatically changes according to the left side type: 
```rust
1 + "1" // 2
"1" + 1 // 11
```

You can escape double quote character within a string:
```rust
"Normal text, \"quoted text\"" // Normal text, "quoted text"
```

## Variables
Kahroba is a dynamically typed language (like Python), the interpreter assigns variables a type at runtime based on the variable's value at the time.

```rust
name = "Kahroba"
version = 0.1
a = 1 + 2
a = "text"
```
## Arrays
You simply define arrays by [] and you can store different types in an array.
```rust
nums = [1,2,3,4]
everything = [1,"kahroba",0.1]
```
Accessing array elements is as easy as below:
```rust
nums[0] // 1
everything[1] // "kahroba"
```
There is also a built-in `len()` function in Kahroba:
```rust
a = [1,2,3,4,5]
println(len(a)) // 5
```
## Maps
Maps in Kahroba are data types for storing key-value pairs of values, they can support any kind of data as their keys or their values:
```rust
data = {
  "name":"Kahroba",
  "version":0.1,
  4: 5
}
println(data["name"])
```
Output:
```
Kahroba
```
**tip:** function `len()` will work on maps and returns the number of pairs.

## Boolean
```rust
a = true
b = false
!a // false
!b // true
a == b // false
a != b // true
```

## Print on terminal
There are two built-in functions `print()` & `println()` for printing values on the terminal.
There would be a printed space between printed arguments and a line break at the end of `println()` output.
```rust
println("Hello world!")
print("Kahroba ")
print("Lang ")
println("version:", 0.1)
```
Output:
```
Hello world!
Kahroba Lang version 0.1
```

## Getting inputs
You can get inputs from users by built-in `input()` function:
```rust
name = input("What is your name:")
print("Hello, ",name)
```
