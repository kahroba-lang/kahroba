## معرفی زبان کهربا
زبان **کهربا** یک پروژه آموزشی جهت نمایش نحوه کار یک زبان برنامه نویسی است
هدف این زبان سادگی و انعطاف پذیری و راحتی استفاده است

[![Documents](https://badgen.net/badge/Documents/passing/blue?icon=wiki)](https://kahroba-lang.github.io/docs/latest)
[![Version](https://badgen.net/badge/Ver/0.1/orange?icon=github)](https://kahroba-lang.github.io/docs/0.1)
[![status](https://badgen.net/badge/Status/Released/green?icon=now)](https://kahroba-lang.github.io/docs/0.1)


## نحوه استفاده از زبان کهربا 
برای مستندات کامل تر لطفا داکیومنت کهربا را از [اینجا](https://kahroba-lang.github.io/docs/latest) بخوانید
.
برای اجرای برنامه های به زبان کهربا لازم است برنامه خود را در فایلی با پسوند .krb یا .kahroba بسازید (مانند main.krb) و از خط فرمان برنامه را اجرا کنید:  
```bash
./kahroba main.krb    // linux
kahroba.exe main.krb   // windows
kahroba_mac main.krb   // mac
```
میتوانید از لینک های زیر فایل های اجرایی مربوط به سیستم عامل خود را دانلود کنید :

[لینوکس](https://github.com/kahroba-lang/kahroba/releases/download/0.1/kahroba) \
[ویندوز](https://github.com/kahroba-lang/kahroba/releases/download/0.1/kahroba.exe) \
[مک](https://github.com/kahroba-lang/kahroba/releases/download/0.1/kahroba_mac)

همچنین میتوانید کهربا را از سورس کد بیلد بگیرید(به گیت و نسخه 1.19 گولنگ نیاز خواهید داشت) :
```
$ git clone https://github.com/kahroba-lang/kahroba.git
$ cd kahroba
$ go build
```

## کامنت

کامنت ها در زبان کهربا با دو اسلش پشت هم شروع میشوند و خطی که کامنت در نظر گرفته شود پردازش نمی شود

    // This is my first program in Kahroba programming language, Let's Rock!

## تعریف رشته ها
رشته ها در زبان کهربا بین دو کوتیشن قرار میگیرند
```rust
"Hello World!"
```
دو رشته را توسط عملگر جمع میتوانید به هم متصل کنید
```rust
"Hello " + "World!" // Hello World
```
در زبان کهربا میتوانید رشته ها را با اعداد جمع کنید  
اگر عملوند سمت چپ رشته باشد نتیجه رشته خواهد شد  
اگر عملوند سمت چپ عدد باشد نتیجه عدد خواهد شد  
مثال:
```rust
1 + "1" // 2
"1" + 1 // 11
```

می‌توان از کاراکترهای کنترلی درون رشته ها استفاده کرد  
مثال:  

```rust
str = "Hello\nWorld"
str2 = "Hello\tWorld"
println(str)
println(str2)
```

خروجی:   
```
Hello
World
Hello	World
```


## تعریف متغیر
متغیرها در زبان کهربا تایپ ندارند و نحوه تعریف کردن اونها به شکل زیر است
```rust
name = "Kahroba"
version = 0.1
a = 1 + 2
a = "text"
```
## آرایه
آرایه ها در زبان کهربا بسیار انعطاف پذیرند و میتوانید در آرایه دیتا تایپ های مختلفی را ذخیره کنید.
```rust
nums = [1,2,3,4]
everything = [1,"kahroba",0.1]
```
برای دسترسی به یک عنصر از آرایه به شکل زیر عمل میکنیم:
```rust
nums[0] // 1
everything[1] // "kahroba"
```
## مپ
مثل آرایه ، مپ هم در زبان کهربا از انعطاف بالایی برخوردار است و میتوانید هر نوع داده ای را درون مپ قرار دهید.
```rust
data = {"name":"Kahroba","version":0.1}
println(data["name"])
```
خروجی:
```
Kahroba
```
## len

به وسیله فانکشن len می‌توانیم طول رشته، آرایه و مپ را بدست بیاوریم.

```rust
a = [1,2,3,4,5]
b = {"name":"kahroba"}
c = "Hello World"

println(len(a)) // 5
println(len(b)) // 1
println(len(c)) // 11
```

## boolean
زبان کهربا از نوع boolean پشتیبانی میکند
```rust
a = true
b = false
!a // false
!b // true
a == b // false
a != b // true
true or false // true
true and false // false
```
## چاپ در خروجی
به وسیله دستور print یا println میتوان عملیات چاپ در خروجی را انجام داد
دستورات چاپ میتوانند چندین ورودی داشته باشند:
```rust
println("سلام دنیا!")
print("زبان ")
print("کهربا ")
println("version:",0.1)
```
خروجی:
```
سلام دنیا!
زبان کهربا version 0.1
```
## گرفتن ورودی
به وسیله دستور input میتوان یک رشته را از کاربر به عنوان ورودی دریافت کرد:
```rust
name = input("What is your name:")
print("Hello, ",name)
```
## تعریف فانکشن
فانکشن ها در زبان کهربا به وسیله کلمه کلیدی fn تعریف میشوند.
فانکشن های میتوانند مقداری باز گردانند یا باز نگردانند.
بصورت پیش فرض آخرین دستور یک فانکشن برگردانده میشود و استفاده از کلمه return اختیاری است
```rust
fn sum(a,b) {
    a+b
}
```
توابع میتوانند بصورت بازگشتی فراخوانی شوند. پیاده سازی مثال کلاسیک فاکتوریل:
```rust
fn f(n) {
    if n <= 1 { 1 }
    n * f(n-1)
}

println(f(5)) // 120
```
ورودی فانکشن میتواند از هر نوعی باشد حتی یک فانکشن دیگر:

    fn getName() {
        "Kahroba"
    }
    fn hello(name) {
        println("Hello ",name)
    }
    hello(getName())

خروجی:

    Hello Kahroba

## swap
توسط این فانکشن میتوانید مقدار دو متغیر را باهم عوض کنید
```rust
a = 5
b = 10
println(a,b)
swap(a,b)
println(a,b)
```
خروجی
```bash
5 10
10 5
```
# کنترل جریان برنامه
## دستورات شرطی
به وسیله دستور if میتوان از دستورات شرطی استفاده کرد
```rust
if a + b > c {
    print("OK")
}
```
همچین میتوان از دستور else برای زمان عدم صحت شرط استفاده کرد
```rust
if a + b > c {
    print("OK")
} else {
    print("Not OK")
}
```
## حلقه تکرار
برای استفاده از حلقه در زبان کهربا از دستور for به شکل زیر استفاده میشود
```rust
for i in 1..5 {
    println(i)
}
```
خروجی:

    1
    2
    3
    4
    5

میتوانید تعداد گام های حلقه را به این شکل مشخص کنید:
```rust
for i in 1..5:2 {
    println(i)
}
```
خروجی:

    1
    3
    5

به وسیله حلقه for میتوانید به روی رشته ها، آرایه ها و مپ ها پیمایش انجام دهید
### پیمایش رشته
```rust
for s in "Hello World" {
    print(s," ")
}
```
خروجی
```bash
H e l l o  W o r l d
```
### پیمایش آرایه
```rust
arr = ["Kahroba","version",0.1]
for v in arr {
    print(v)
}

for i,v in arr {
    println(i,":",v)
}
```

خروجی
```bash
Kahroba version 0.1

0:Kahroba
1:version
2:0.1
```
### پیمایش مپ
```rust
data = {"name":"Kahroba","version":0.1}
for v in data {
    println(v)
}

for k,v in data {
    println(k,":",v)
}
```
خروجی
```bash
Kahroba
0.1

name : Kahroba
version : 0.1
```

### فراخوانی فایل های کهربا خارجی
با استفاده از زبان کهربا همچنین می توانید از سایر فایل های .krb که به زبان کهربا نوشته شده است استفاده کنید. 

```rust 
import("helper.krb")
```

تمامی فانکشن ها و متغیرهای تعریف شده در فایل helper.krb در اسکوپ جاری برنامه قابل دسترس خواند بود.


## پیاده سازی الگوریتم quicksort توسط کهربا
```rust
fn qsort(arr) {
    sort(arr,0,len(arr)-1)
}

fn sort(arr,l,r) {
    if l < r {
        q = partition(arr,l,r)
        sort(arr,l,q-1)
        sort(arr,q+1,r)
    }
}

fn partition(arr,l,r)  {
    i = l
    for j in l..r {
        if a[j] < a[r] {
            swap(a[i],a[j])
            i = i + 1
        }
    }
    swap(a[i],a[r])
    i
}


a = [5,1,2,4,3,9,8,7,6,0]

println(a)
qsort(a)
println(a)
```
خروجی

```bash
[5 1 2 4 3 9 8 7 6 0]
[0 1 2 3 4 5 6 7 8 9]
```
> [مثال های بیشتر](/examples/)
