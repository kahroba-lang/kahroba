/*
مقدمه
به هیجان انگیز ترین قسمت پیاده سازی زبان خوش اومدید
تمام جادو اینجا اتفاق میفتد و ما باید در برخورد با هر نود بتونیم دستورات
مربوط به آن نود در زبان مبدا را مشخص کنیم که اینجا زبان مبدا گو هست
همه نود ها اینترفیس
Node
را پیاده سازی کرده اند پس به محض ورود یک نود به چنل نود ها از مرحله قبل
بلافاصله اقدام به اجرای متد
Eval
مربوط به آن نود میکنیم و دستوارت مربوط به آن نود رو اجرا میکنیم
*/
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

/*
مهمترین قدم برای اجرای دستورات مشخص کردن اسکوپ دستورات است
باید بدانیم متغیری که در حال استفاده از آن هستیم درون چه اسکوپی هست و چگونه به آن دسترسی داشته باشیم.
ایده کار ساده است یک استراکت که دو مپ درون خودش نگهداری میکند یکی برای متغیرها و دیگری برای
نگهداری فانکشن ها. البته میشود از یک مپ هم استفاده کرد ولی در آن صورت نمیتوان فانکشن همنام با یک متغیر درون یک اسکوپ داشته باشیم
کاری که گو اجازه آن را نمی‌دهد ولی ما در کهربا میتونیم یک فانکشن هم نام با یک متغیر داشته باشیم :) ‌ـ
*/
type Scope struct {
	variables map[string]any
	functions map[string]any
	parent    *Scope
}

/*
یک سازنده که مقدار دهی اولیه را انجام میدهد
همانطور که می بینید یک متغیر به نام پرنت داریم که کارش این است که به
اسکوپ بالاسری خودش اشاره کند
اگر اسکوپ ما پرنت نداشته باشد یعنی اسکوپ گلوبال است
*/
func NewScope(scope *Scope) *Scope {
	return &Scope{
		variables: make(map[string]any),
		functions: make(map[string]any),
		parent:    scope,
	}
}

/*
یک متغیر را در لیست متغیر ها ثبت میکند
*/
func (e Scope) SetVariable(k Node, v any) {
	switch node := k.(type) {
	case Identifier:
		e.variables[node.Token.Value] = v
	case ArrayMapIndex:
		name := node.Subject.(Identifier).Token.Value
		index := node.Index.Eval(&e)
		t, _ := e.GetVariable(name)
		switch arrMap := t.(type) {
		case []any:
			arrMap[index.(int)] = v
		case map[any]any:
			arrMap[index] = v
		}
		e.variables[name] = t
	}
}

/*
یک متغیر را از لیست متغیر ها فراخوانی میکند
اگر متغیر در اسکوپ فعلی پیدا نشد این کار را تا رسیدن به اسکوپ گلوبال بصورت بازگشتی ادامه میدهد
*/
func (e Scope) GetVariable(s string) (any, bool) {
	v, ok := e.variables[s]
	if !ok && e.parent != nil {
		v, ok = e.parent.GetVariable(s)
	}
	return v, ok
}

/*
یک فانکشن را در لیست فانکشن ها ثبت میکند
*/
func (e Scope) SetFunction(s string, v any) {
	e.functions[s] = v
}

/*
یک فانکشن را از لیست فانکشن ها فراخوانی میکند
اگر فانکشن در اسکوپ فعلی پیدا نشد این کار را تا رسیدن به اسکوپ گلوبال بصورت بازگشتی ادامه میدهد
*/
func (e Scope) GetFunction(s string) (any, bool) {
	v, ok := e.functions[s]
	if !ok && e.parent != nil {
		v, ok = e.parent.GetFunction(s)
	}
	return v, ok
}

/*
این فانکشن روی چنل نود ها پیمایش میکند و متد
Eval
مربوط به هر نود را اجرا میکند
*/
func Eval(nodes chan Node, scope *Scope) any {
	var r any
	for node := range nodes {
		r = node.Eval(scope)
	}
	return r
}

/*
مثل پارسر اینجا هم با نود ساده تر شروع میکنیم
برای اجرای نودهای زیر کافیست مقدار ذخیره شده درون نود را برگردانیم
String,Bool,Int,Float,Return
شاید اینجا بپرسید ما گفتیم هر نود مسئول اجرای یک سری دستور است ولی اینجا ما چیزی اجرا نمیکنیم
فقط یک مقدار را ریترن میکنیم
اما در ادامه خواهید دید مقدار بازگشتی از این نود ها توسط نودهای دیگر مورد استفاده قرار میگیرند
*/
func (n String) Eval(scope *Scope) any {
	return n.Value
}

func (n Bool) Eval(scope *Scope) any {
	return n.Value
}

func (n Int) Eval(scope *Scope) any {
	return n.Value
}

func (n Float) Eval(scope *Scope) any {
	return n.Value
}

func (n Return) Eval(scope *Scope) any {
	return n.Value.Eval(scope)
}

/*
برای اجرای آیدنتیفایرها کافیست بدانیم
که این آیدنتیفایر به یک فانکشن اشاره میکند یا یک متغیر
با فلگی که در مرحله پارس کردن گذاشتیم کار ما اینجا ساده شده است و میتوانیم متد مناسب از روی اسکوپ را صدا بزنیم
*/
func (n Identifier) Eval(scope *Scope) any {
	if n.Function {
		v, _ := scope.GetFunction(n.Token.Value)
		return v
	}
	v, _ := scope.GetVariable(n.Token.Value)
	return v
}

/*
برای مقدار دهی متغیر کافیست مقدار متغیر را محاسبه کرده و درون اسکوپ مقدار آن را ذخیره کنیم
*/
func (n Variable) Eval(scope *Scope) any {
	v := n.Value.Eval(scope)
	scope.SetVariable(n.Name, v)
	return v
}

/*
برای اجرای پیشوند ها مثل
-5
!true
ابتدا مقدار ذخیره شده درون نود را محاسبه میکنیم سپس با توجه به نوع پیشوند
آن را منفی یا
not
میکنیم
*/
func (n UnaryOperator) Eval(scope *Scope) any {
	v := n.Exp.Eval(scope)
	return evalPrefix(n.Token.Value, v)
}

func evalPrefix(prefix string, v any) any {
	switch v.(type) {
	case bool:
		if prefix == "!" {
			return !v.(bool)
		}
	case int:
		return v.(int) * -1
	case float64:
		return v.(float64) * -1
	}
	return nil
}

/*
برای محاسبه عملیات باینری باید ابتدا سمت چپ و راست معادله را حساب کرد
سپس بر اساس عملگر تصمیم گرفت در حالت های مختلف باید چه کاری کرد
از آنجایی که در زبان کهربا دستورات زیر خروجی های متفاوتی را در بر میگیرد
پس لازم است همه حالت های مختلف را در نظر گرفت و برایشان محاسبه انجام داد
مثال:
1 + "1" => 2
"1" + 1 => 11
در ادامه یک سری فانکشن کمکی نوشتیم که حالت های مختلف را محاسبه میکنند
تقریبا همه آنها شبیه به هم هستند ولی از آنجایی که گو یک زبان
تایپ سیف است پس نیاز داریم که حالت های مختلف را با نوع داده مربوط به خود محاسبه کنیم
هرچند در ظاهر همه آنها شبیه به هم باشند
*/
func (n BinaryOperator) Eval(scope *Scope) any {
	l := n.Left.Eval(scope)
	r := n.Right.Eval(scope)
	operator := n.Token.Value

	if l == nil {
		return r
	}
	if r == nil {
		return l
	}

	switch l.(type) {

	case int:

		switch r.(type) {
		case int:
			return evalIntInt(l, r, operator)
		case string:
			return evalIntString(l, r, operator)
		case float64:
			return evalIntFloat(l, r, operator)
		default:
			return nil
		}

	case float64:

		switch r.(type) {
		case int:
			return evalIntFloat(r, l, operator)
		case string:
			return evalFloatString(l, r, operator)
		case float64:
			return evalFloatFloat(l, r, operator)
		default:
			return nil
		}

	case string:
		switch r.(type) {
		case string:
			return evalStringString(l, r, operator)
		case int:
			return evalStringInt(l, r, operator)
		case float64:
			return evalStringFloat(l, r, operator)
		default:
			return nil
		}

	case bool:
		switch r.(type) {
		case bool:
			return evalBoolBool(r, l, operator)
		default:
			return nil
		}

	default:
		return nil
	}
}

func evalIntInt(ll, rr any, operator string) any {
	l := ll.(int)
	r := rr.(int)
	switch operator {
	case ">":
		return l > r
	case ">=":
		return l >= r
	case "<":
		return l < r
	case "<=":
		return l <= r
	case "==":
		return l == r
	case "!=":
		return l != r
	case "+":
		return l + r
	case "-":
		return l - r
	case "*":
		return l * r
	case "/":
		return l / r
	default:
		return nil
	}
}

func evalFloatFloat(ll, rr any, operator string) any {
	l := ll.(float64)
	r := rr.(float64)
	switch operator {
	case ">":
		return l > r
	case ">=":
		return l >= r
	case "<":
		return l < r
	case "<=":
		return l <= r
	case "==":
		return l == r
	case "!=":
		return l != r
	case "+":
		return l + r
	case "-":
		return l - r
	case "*":
		return l * r
	case "/":
		return l / r
	default:
		return nil
	}
}

func evalIntString(ll, rr any, operator string) int {
	l := ll.(int)
	r, _ := strconv.Atoi(rr.(string))
	switch operator {
	case "+":
		return l + r
	case "-":
		return l - r
	case "*":
		return l * r
	case "/":
		return l / r
	default:
		return 0
	}
}

func evalFloatString(ll, rr any, operator string) float64 {
	l := ll.(float64)
	r, _ := strconv.ParseFloat(rr.(string), 64)
	switch operator {
	case "+":
		return l + r
	case "-":
		return l - r
	case "*":
		return l * r
	case "/":
		return l / r
	default:
		return 0
	}
}

func evalIntFloat(ll, rr any, operator string) any {
	l := float64(ll.(int))
	r := rr.(float64)
	switch operator {
	case ">":
		return l > r
	case ">=":
		return l >= r
	case "<":
		return l < r
	case "<=":
		return l <= r
	case "==":
		return l == r
	case "!=":
		return l != r
	case "+":
		return l + r
	case "-":
		return l - r
	case "*":
		return l * r
	case "/":
		return l / r
	default:
		return nil
	}
}

func evalStringString(ll, rr any, operator string) string {
	l := ll.(string)
	r := rr.(string)
	switch operator {
	case "+":
		return l + r
	default:
		return ""
	}
}

func evalStringInt(ll, rr any, operator string) string {
	l := ll.(string)
	r := strconv.Itoa(rr.(int))
	switch operator {
	case "+":
		return l + r
	default:
		return ""
	}
}

func evalStringFloat(ll, rr any, operator string) string {
	l := ll.(string)
	r := strconv.FormatFloat(rr.(float64), 'f', -1, 64)
	switch operator {
	case "+":
		return l + r
	default:
		return ""
	}
}

func evalBoolBool(ll, rr any, operator string) bool {
	r := rr.(bool)
	l := ll.(bool)
	switch operator {
	case "==":
		return l == r
	case "!=":
		return l != r
	case "or":
		return l || r
	case "and":
		return l && r
	}
	return false
}

/*
اجرای بلاک هم کار ساده ای است
کافیست روی اسلایس نود ها پیمایش کنیم و متد
Eval
هر نود را اجرا کنیم
تنها نکته ای که وجود دارد این است که اگر میان یکی از دستورات به دستور
return یا if
برخورد کردیم مقدار بازگشتی از آنها را برگردانیم و بقیه را اجرا نکنیم
در مورد کلمه کلیدی
return
که علت واضح است اما در مورد دستور
If
نیاز به کمی توضیح داریم
در زبان کهربا آخرین دستور بلاک شرط یا فانکشن برگشت داده میشود و نیاز به تایپ
return
نیست
برای مثال فانکشن های زیر کاملا معتبر هستند
fn sum(a,b) {

	a+b // return a+b

}

fn num() {

	1
	2
	3 // return 3

}

پس اگر وسط دستوارتمان یک شرط داشتیم و نتیجه اجرای آخرین دستور آن
nil
نبود، یعنی مثلا آخرین دستور پرینت نبود چون عبارات غیر محاسبه ای مقدار
nil
برمیگرداند
در این حالت یعنی یک چیزی از دستور شرط ما برگشت داده شده است
پس اگر نوع نود
return یا if
بود همانجا نتیجه را برمیگردانیم و دستورات بعدی را اجرا نمیکنیم
*/
func (n Block) Eval(scope *Scope) any {
	var result any
	for _, stm := range n.Statements {
		if stm == nil {
			continue
		}
		result = stm.Eval(scope)
		if result != nil {
			switch stm.(type) {
			case If, Return:
				return result
			}
		}
	}
	return result
}

/*
دستورات شرطی هم به راحتی قابل محاسبه هستند
شرط را محاسبه میکنیم اگر نتیجه اجرای شرط درست بود بلاک کد
True
در غیر اینصورت بلاک کد
Else
را در صورت وجود اجرا میکنیم
*/
func (n If) Eval(scope *Scope) any {
	condition := n.Condition.Eval(scope).(bool)
	if condition {
		return n.True.Eval(scope)
	} else if n.Else != nil {
		return n.Else.Eval(scope)
	}
	return nil
}

/*
برای محاسبه آرایه/مپ نیاز است که ابتدا تمام عضو های آن را محاسبه کنیم
چون ممکن است عضوهای خودشان عبارت یا فانکشن باشند
مثال
a = [1,5+3,4/3,hello()]
*/
func (n Array) Eval(scope *Scope) any {
	ret := []any{}
	for _, node := range n.Nodes {
		ret = append(ret, node.Eval(scope))
	}
	return ret
}
func (n Map) Eval(scope *Scope) any {
	m := map[any]any{}
	for k, v := range n.Nodes {
		m[k.Eval(scope)] = v.Eval(scope)
	}
	return m
}

/*
برای محاسبه اندیس یک آرایه یا مپ کافیست متغیر آرایه یا مپ درون نود را محاسبه کنیم
سپس ایندکس را محاسبه کنیم و ایندکس آن متغیر را برگردانیم
در مورد آرایه همیشه اندیس عدد صحیح است ولی در مورد مپ، اندیس میتواند هرچیزی باشد
*/
func (n ArrayMapIndex) Eval(scope *Scope) any {
	index := n.Index.Eval(scope)
	arrMap := n.Subject.Eval(scope)
	switch arrMap.(type) {
	case map[any]any:
		return arrMap.(map[any]any)[index]
	default:
		return arrMap.([]any)[index.(int)]
	}
}

/*
محاسبه دستور پرینت هم ساده است ابتدا باید آرگومان ها را محاسبه کنیم و بر اساس فلگی که قبلا
ست کردیم یکی از فانکشن های
print یا println
را با آرگومان هایی که محاسبه کردیم اجرا کنیم
*/
func (n Print) Eval(scope *Scope) any {
	args := evalExpressions(n.Args, scope)
	if n.NewLine {
		fmt.Println(args...)
	} else {
		fmt.Print(args...)
	}
	return nil
}

/*
این فانکشن کمکی هر نود را اجرا میکند و نتیجه را درون یک اسلایس قرار میدهد
*/
func evalExpressions(exps []Node, scope *Scope) []any {
	res := []any{}
	for _, exp := range exps {
		r := exp.Eval(scope)
		res = append(res, r)
	}
	return res
}

/*
برای محاسبه یک فانکشن کافیست که اسکوپ را درون استراکت نود اینجکت کنیم و در لیست
فانکشن ها آن را ثبت کنیم
*/
func (n Function) Eval(scope *Scope) any {
	n.Scope = scope
	scope.SetFunction(n.Name, n)
	return n
}

/*
برای محاسبه فانکشن کال باید ابتدا خود فانکشن را محاسبه کنیم
سپس آرگومان های فانکشن را محاسبه میکنیم
سپس توسط فانکشن
applyfunction
فانکشنن را اجرا میکنیم
*/
func (n FunctionCall) Eval(scope *Scope) any {
	fn := n.Function.Eval(scope).(Function)
	args := evalExpressions(n.Args, scope)
	return applyfunction(fn, args, true)
}

/*
قبل از اجرای فانکشن نیاز هست که آرگومان ها را به پارامتر ها وصل کنیم
در تعریف فانکشن پارامتر داریم که مقدار ندارد
در صدا زدن فانکشن آرگومان داریم که مقدار دارد ولی نام متغیر ندارد
از طریق فانکشن زیر این دو را بهم متصل میکنیم
مثال:
fn sum(a,b) {...}
sum(5,7)
باید در مثال بالا متغیرها را مقدار دهی کنیم
a=5
b=7
در اینجا شاید علت استفاده از اسکوپ برایتان شفاف شود
تا بحال فقط یک متغیر از نوع اسکوپ داشتیم که بین نود های مختلف پاس میدادیم
ولی برای اجرای یک فانکشن ما نیاز به اسکوپ مخصوص به خودش رو داریم
در اینجا ما اسکوپ جدیدی میسازیم و اسکوپ قبلی را پرنت اسکوپ فعلی قرار میدیم
به این شکل فانکشن ما دارای اسکوپ خودش خواهد بود چون متغیر هایی که به عنوان پارامتر
استفاده شده یا متغیرهای درون یک فانکشن فقط مختص خود آن فانکشن هستند و نباید با اسکوپ های دیگر تداخل داشته باشد
همانطور که می‌بینید ما با یک فلگ مشخص میکنیم که آیا اسکوپ جدید بسازیم یا نه
علت استفاده از این فلگ این است که ما از این فانکشن درون حلقه
for
هم استفاده میکنیم و آنجا نیازی نیست که در هر بار اجرای حلقه اسکوپ جدید استفاده شود
وقتی به توضیح حلقه رسیدیم این موضوع را بیشتر شفاف میکنم
*/
func argsToScope(fn Function, args []any, new bool) *Scope {
	scope := fn.Scope
	if new {
		scope = NewScope(fn.Scope)
	}
	for i, param := range fn.Params {
		scope.SetVariable(*param, args[i])
	}
	return scope
}

/*
بعد از اینکه فانکشن محاسبه شد مقدار آرگومان ها را به پارامتر ها وصل میکنیم سپس نوبت اجرای دستورات فانکشن است
که قبلا در اجرای بلاک با نحوه اجرای آن آشنا شدید
*/
func applyfunction(fn Function, args []any, new bool) any {
	newScope := argsToScope(fn, args, new)
	return fn.Body.Eval(newScope)
}

/*
برای اجرای یک حلقه
for
ما با آن مثل یک فانکشن رفتار میکنیم و در هر بار اجرای حلقه
به شمارنده یک مقدار جدید میدهیم و انگار فانکشنی را با یک مقدار جدید صدا زدیم
تنها تفاوتش با فانکشن کال معمولی این است که اسکوپ حلقه نباید در هربار اجرا تغییر کند
فرض کنید میخواهیم شکل زیر را چاپ کنیم
*
**
***
****
حلقه فور باید به این شکل باشد

for i in 1..4 {

	str = ""
	for j in 1..i {
		str = str + "*"
	}
	println(str)

}

اگر دقت کنید ما درون حلقه داخلی داریم به مقدار قبلی رشته یک ستاره اضافه میکنیم
اگر اسکوپ حلقه در هر بار اجرا تغییر کند مقدار
str
همیشه خالی خواهد بود چون به مقدار قبلی رشته دسترسی ندارد
برای اینکه این مشکل را حل کنیم فلگ
new
را برای اجرای دستورات حلقه
for
بصورت
false
ارسال میکنیم
از آنجایی که حلقه ما میتواند روی آرایه ، مپ و رشته پیمایش انجام دهد باید برای هر کدام از این
حالت ها متغیر سابجکت را کست کنیم و دستورات را بوسیله فانکشن
applyfunciton
اجرا کنیم
*/
func (n For) Eval(scope *Scope) any {
	subject := n.Subject.Eval(scope)
	fn := Function{
		Body:   n.Block,
		Params: []*Identifier{n.Key},
		Scope:  scope,
	}

	if n.Value != nil {
		fn.Params = append(fn.Params, n.Value)
	}

	switch subject.(type) {
	case map[any]any:
		for k, v := range subject.(map[any]any) {
			args := []any{k}
			if n.Value != nil {
				args = []any{k, v}
			}
			applyfunction(fn, args, false)
		}
	case string:
		for _, v := range subject.(string) {
			args := []any{string(v)}
			applyfunction(fn, args, false)
		}
	case []any:
		for k, v := range subject.([]any) {
			args := []any{v}
			if n.Value != nil {
				args = []any{k, v}
			}
			applyfunction(fn, args, false)
		}
	}

	return nil
}

/*
رنج یک اسلایس از اعداد صحیح میسازد
قسمت گام اختیاریست و اگر گام ست نشده باشد گام ۱ در نظر گرفته میشود
همچنین حالت رنج صعودی و نزولی هم در نظر گرفته شده
مثال
1..5    // 1 2 3 4 5
1..5:2  // 1 3 5
5..1    // 5 4 3 2 1
5..1:2  // 5 3 1
*/
func (n Range) Eval(scope *Scope) any {
	from := n.From.Eval(scope).(int)
	to := n.To.Eval(scope).(int)
	step := 1
	if n.Step != nil {
		step = n.Step.Eval(scope).(int)
	}

	ret := []any{}
	if from < to {
		for i := from; i <= to; i += step {
			ret = append(ret, i)
		}
		return ret
	}
	if step > 0 {
		step *= -1
	}
	for i := from; i >= to; i += step {
		ret = append(ret, i)
	}
	return ret

}

/*
swap
برای جابجا کردن مقدار دو متغیر استفاده میشود
برای اینکه مقدار متغیر اول ما از بین نرود از یک متغیر میانی برای این کار استفاده میکنیم
*/
func (n Swap) Eval(scope *Scope) any {
	t := n.A.Eval(scope)
	scope.SetVariable(n.A, n.B.Eval(scope))
	scope.SetVariable(n.B, t)
	return nil
}

/*
import
برای import کردن یک فایل استفاده می شود
*/
func (n Import) Eval(scope *Scope) any {
	t := n.Filename.Eval(scope).(string)

	input, err := os.ReadFile(t)
	if err != nil {
		log.Fatal(err)
	}
	l := NewLexer(string(input))
	p := NewParser(l.tokens)

	Eval(p.nodes, scope)

	return nil
}

/*
input
برای دریافت ورودی از کاربر استفاده میشود
یک متغیر را به عنوان ورودی میگیرد و مقدار آن را برابر با رشته ورودی قرار میدهد
*/
func (n Input) Eval(scope *Scope) any {
	reader := bufio.NewReader(os.Stdin)
	prompt := n.Promp.Eval(scope)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return text
}

/*
برای اینکه بتوانیم طول آرایه یا مپ را محاسبه کنیم ابتدا باید آن را به دیتا تایپ
اصلی کست کنیم تا زبان گو بتواند فانکشن
len
را روی آن متغیر اجرا کند
*/
func (n Len) Eval(scope *Scope) any {
	v := n.ArrMap.Eval(scope)
	switch t := v.(type) {
	case string:
		return len(t)
	case map[any]any:
		return len(t)
	case []any:
		return len(t)
	}
	return nil
}

/*
تبریک میگم ساخت زبان به پایان رسید
دو فایل بعدی فقط کاغذ بازی اداریه و میخوایم تمام چیزهایی که نوشتیم رو کنار هم بگذاریم
راه طولانی ولی لذت بخش بود
چیزی تا انتها نمونده :)
*/
