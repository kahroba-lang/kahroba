/*
مقدمه
بعد از تقسیم فایل به توکن ها ، حالا نوبت آن است که توکن های مرتبط به هم را در کنار هم قرار دهیم
و یک دستور معنی دار را شکل دهیم برای مثال اگر به توکن
if
برخورد کردیم توکن های بعد از آن را به شکلی تفسیر میکنیم که یک دستور شرطی کامل را شکل دهد
برای مثال بعد از توکن
if
تمام توکن ها را به عنوان شرط اجرای حلقه در نظر میگیریم تا به توکن
{
برسیم که به معنی شروع دستورات مربوط به حلقه شرطمان است
و همین روش را برای باقی توکن ها انجام میدهیم تا از ترکیب توکن ها
یک عبارت یا دستور کامل را شکل دهیم. هدف اجرای پارسر این است که در انتها یک درخت از دستورات را شکل دهیم
به این درخت
Abstract Syntax Tree
یا در اختصار
AST
گفته میشود.
در مرحله بعد کار ما پیمایش درخت و اجرای دستورات مربوط به هر نود درخت میباشد.
دستور زیر را در نظر بگیرید
1 + 2 * 3
در نهایت تبدیل به درخت زیر خواهد شد
BinaryOperator
├── 1
├── +
└── BinaryOperator

	    ├── 2
		├── *
	    └── 3
*/
package main

import (
	"strconv"
)

/*
در اینجا ما لیستی از عملگرها را به ترتیب اولویت آنها مرتب کرده ایم
دنیا جای قشنگ تری بود اگر از سمت چپ شروع میکردیم و عملگرها را اعمال میکردیم
ولی در دنیای واقعی نتیجه عبارت زیر ۱۶ نیست زیرا عملگر ضرب اولویت بالاتری از جمع دارد
3 + 5 * 2
پس باید در زمان محاسبه یک عبارت بدانیم عملگری که
در حال پردازش آن هستیم اولویت بالاتری از باقی عملگر ها دارد یا کمتر
*/
const (
	LOWEST_PRIORITY = iota + 1
	EQUALS          // ==
	LESSGREATER     // > < >= <=
	SUM             // +
	PRODUCT         // *
	PREFIX          // -5 or !true
	CALL            // myfunc()
	INDEX           // a[2]
	RANGE           // 1..10
)

/*
اولویت ها را درون یک مپ قرار میدهیم تا به سادگی به اولویت آن دست پیدا کنیم
*/
var priorities = map[Type]int{
	EQ:       EQUALS,
	EQEQ:     EQUALS,
	NEQ:      EQUALS,
	LESSER:   LESSGREATER,
	GREATER:  LESSGREATER,
	LEQ:      LESSGREATER,
	GEQ:      LESSGREATER,
	PLUS:     SUM,
	MINUS:    SUM,
	SLASH:    PRODUCT,
	STAR:     PRODUCT,
	LPARENT:  CALL,
	LBRACKET: INDEX,
	DOTDOT:   RANGE,
}

/*
تمام توکن های که پارس (پردازش) میشوند در نهایت تبدیل به یک نود میشوند
اینترفیس نود فقط یک متد دارد که برای اجرای دستورات مربوط به هر نود به کار میرود
در فایل بعدی که مربوط به اجرای دستورات است، برای هر نود دستورات مربوط به آن نود را پیاده سازی میکنیم
*/
type Node interface {
	Eval(scope *Scope) any
}

/*
ایده کلی پارسر به این صورت است که به هر توکنی که رسیدیم یک متد متناسب با آن توکن را کال می‌کنیم
دقیقا مانند کاری که در لکسر انجام دادیم. با این تفاوت که در لکسر ما کاراکتر به کاراکتر بررسی میکردیم
در اینجا ما توکن به توکن بررسی میکنیم
تمام عملیاتی که بر روی توکن ها انجام میدهیم و پارس کردن آنها از دو حالت خارج نیست
یا توکن ها به تنهایی دارای معنی هستند مثلا توکن
for
به تنهایی نشان دهنده شروع یک حلقه است و میتوانیم تمام توکن های بعد از آن را به راحتی تفسیر کنیم
ولی یک عملگر مانند عملگر منفی به تنهایی قابل تفسیر نیست
اگر به عملگر منفی رسیدیم به این معنی است که میخواهیم عدد بعد از خودش را منفی کنیم؟
یعنی قصدمان تفسیر عدد 5- است؟ یا نه این منفی در یک عبارت ریاضی به کار برده شده
و باید از عدد قبل از خود کم شود. مثلا
17-5
بخاطر همین ما پارس کردن توکن ها را به دو بخش تقسیم میکنیم . توکن هایی که به تنهایی قابل تفسیر هستند
و توکن هایی که به تنهایی قابل تفسیر نیستند
برای راحتی کار دو نوع فانکشن تعریف میکنیم برای هر کدام از انواع عملیات باینری و یونری
عملیات یونری چون به تنهایی قابل تفسیر هستند هیچ ورودی ای ندارند
ولی عملیات باینری چون نیاز به یک عملوند دیگر برای معنی دار شدند دارند آن قسمت دوم را بصورت
ورودی به فانکشن پاس میدهیم
*/
type unaryFn func() Node
type binaryFn func(Node) Node

/*
پارسر ما دارای چند متغیر درون خود هست که آنها را باهم بررسی میکنیم
nodes
چنلی که بعد از پارس کردن توکن ها دستور عمل بدست آمده را درون آن قرار میدهیم تا قسمت بعدی یعنی
evaluator
از آن استفاده کند.
دقیقا مانند کاری که در لکسر کردیم و خروجی لکسر توسط پارسر پردازش شد. در اینجا خروجی لکسر توسط
evaluator
محاسبه میشود.
این چنل از نوع
Node
هست. هر عبارت را یک آبجکت از نوع نود در نظر گرفتیم که در فایل بعدی به تفصیل به آن خواهیم پرداخت.

tokens
چنلی از توکن ها که همان خروجی لکسر است که توکن ها درون آن قرار میگیرند
توکن ها را از این چنل استخراج میکنیم و به تفسیر آن میپردازیم

currentToken
توکنی که در حال بررسی آن هستیم

nextToken
توکنی که بعد از توکن فعلی در صف پردازش قرار دارد

unaryOperations و binaryOperations
یک مپ برای وصل کردن هر توکن به متدی که قرار است آن توکن را پردازش کند
*/
type parser struct {
	nodes            chan Node
	tokens           chan Token
	currentToken     *Token
	nextToken        *Token
	unaryOperations  map[Type]unaryFn
	binaryOperations map[Type]binaryFn
}

/*
نکته خاصی در این سازنده وجود ندارد. از طریق این سازنده پارسرمان را ایجاد میکنیم و مقادیر اولیه را درون آن
ست میکنیم سپس استارت پردازش توکن ها را با اجرای متد
parse
بصورت گورتین میزنیم . تا مرحله بعدی یعنی
evaluator
منتظر اتمام پردازش همه توکن ها نماند و به محض پردازش یک نود جدید عملیات را روی آن انجام دهد.
در این سازنده، مپ های متد ها را به هر توکن متصل میکنیم و میگوییم هر توکن با چه متدی پردازش شود.
*/
func NewParser(tokens chan Token) *parser {
	p := &parser{
		tokens:           tokens,
		binaryOperations: make(map[Type]binaryFn),
		nodes:            make(chan Node),
		currentToken:     &Token{},
		nextToken:        &Token{},
	}

	/*
		توکن هایی که به تنهایی قابل تفسیر هستند
	*/
	p.unaryOperations = map[Type]unaryFn{
		IDENT:    p.parseIdentifier,
		STRING:   p.parseString,
		INT:      p.parseInt,
		FLOAT:    p.parseFloat,
		MINUS:    p.parseUnaryOperator,
		NOT:      p.parseUnaryOperator,
		TRUE:     p.parseBool,
		FALSE:    p.parseBool,
		LPARENT:  p.parseGrouped,
		IF:       p.parseIf,
		FN:       p.parseFunction,
		PRINT:    p.parsePrint,
		PRINTLN:  p.parsePrint,
		LBRACKET: p.parseArray,
		LCURLY:   p.parseMap,
		FOR:      p.parseFor,
		RETURN:   p.parseReturn,
		SWAP:     p.parseSwap,
		INPUT:    p.parseInput,
		LEN:      p.parseLen,
		IMPORT:   p.parseImport,
	}

	/*
		عملگرهای ریاضی و منطقی که نیاز به دو عملوند دارند و بصورت باینری پردازش میشوند
	*/
	for _, typ := range []Type{OR, AND, PLUS, MINUS, STAR, SLASH, EQEQ, NEQ, GREATER, GEQ, LESSER, LEQ} {
		p.binaryOperations[typ] = p.parseBinaryOperation
	}

	/*
		توکن هایی که باید بصورت باینری به همراه یک توکن دیگر تفسیر شوند
	*/
	p.binaryOperations[LPARENT] = p.parseFunctionCall
	p.binaryOperations[LBRACKET] = p.parseArrayIndex
	p.binaryOperations[DOTDOT] = p.parseRange
	p.binaryOperations[EQ] = p.parseVariable

	// توکن فعلی و توکن بعدی را مقدار دهی اولیه میکنیم
	*p.currentToken = <-p.tokens
	*p.nextToken = <-p.tokens

	go p.parse()

	return p
}

/**************
* متدهای کمکی *
**************/

/*
این متد اولویت توکن مورد نظر را در جدول اولویت توکن ها پیدا میکند
اگر اولویت توکن مورد نظر را پیدا نکرد کمترین اولویت را برمیگرداند
*/
func (p *parser) getPriority(tkn Type) int {
	if n, ok := priorities[tkn]; ok {
		return n
	}
	return LOWEST_PRIORITY
}

/*
چک میکند که آیا توکن بعدی از نوع مورد نظر ما است یا نه
اگر بود به توکن بعدی میرود
*/
func (p *parser) isNextToken(t Type) bool {
	if p.nextToken.Type == t {
		p.next()
		return true
	}
	return false
}

/*
این متد مقدار توکن بعدی را در توکن فعلی قرار میدهد
و یک توکن جدید از چنل دریافت میکند به عنوان توکن بعدی
*/
func (p *parser) next() {
	p.currentToken = p.nextToken
	t := <-p.tokens
	p.nextToken = &t
}

/*********************
*  پایان متدهای کمکی *
*********************/

/*
در این متد تا وقتی به توکن
EOF
نرسیده ایم توکن فعلی را چک میکنیم و بصورت یک نود به چنل نود ها اضافه میکنیم
این متد توکن فعلی را میخواند و پارس میکند تا به انتهای آن عبارت برسد و پس از تکمیل یک عبارت
به سراغ عبارت بعدی میرود. توکن اولی که این متد چک میکند اولین توکن فایل است ولی ممکن است بار دیگری که این
حلقه اجرا میشود ۵۰ توکن بعد باشد چون هر متدی که کال میشود تعدادی از توکن ها را مصرف میکند تا به یک
عبارت معنی دار برسد
*/
func (p *parser) parse() {
	for p.currentToken.Type != EOF {
		node := p.parseExpression(LOWEST_PRIORITY)
		if node != nil {
			p.nodes <- node
		}
		p.next()
	}
	close(p.nodes)
}

/*
مهمترین متد پارسر این متد است
کاری که این متد انجام میدهد این است که به محض برخورد با یک توکن، متد مرتبط به آن توکن را فراخوانی میکند
الگوریتم این متد بسیار ساده اما خیلی کارآمد و قدرتمند است و به این شکل عمل میکند
از آنجایی که ما هم توکن های باینری داریم و هم توکن های که
یونری، نیاز هست که هر دو حالت در نظر گرفته شود
در ابتدا توکن مورد نظر در فهرست توکن های یونری جستجو میشود و نتیجه اجرای متد مورد نظر در متغیر
left
ریخته میشود. سپس به سراغ توکن بعدی می‌رویم، اگر توکن در فهرست توکن های باینری پیدا شد
پس متغیر
left
سمت چپ معادله است و سمت راست معادله میشود توکن بعدی که در صف پردازش قرار دارد.
اما اگر بعد از پردازش
left
توکن بعدی یک توکن باینری نبود همان متغیر
left
برگشت داده میشود!
شاید کمی گیج کننده باشد پس اجازه دهیم با یک مثال موضوع را مشخص کنیم
فرض کنید در حال پردازش این عبارت هستیم

return 5

به محض برخورد با توکن اول که توکن ریترن است آن را درون توکن های یونری سرچ میکنیم
می‌بینیم که در لیست توکن های یونری توکن ریترن پیدا میشود و شروع به پردازش این عبارت میکنیم و متد
parseReturn
توکن ریترن و توکن بعد از خودش یعنی 5 را میخواند و پردازش این عبارت به اتمام میرسند.
در واقع متغیر لفت در برگیرنده پردازش توکن ریترن و عدد 5 خواهد بود
بعد از اجرای متد
parseReturn
دوباره کنترل برنامه به متد اصلی که این متد را کال کرده برمیگردیم  یعنی متد
parseExpression
بعد از پر شدن متغیر لفت به سراغ باقی دستوات این فانکشن میرویم از آنجایی که توکنی وجود ندارد
پس متغیر لفت به عنوان یک عبارت کامل برگشت داده میشود.
حال این عبارت را در نظر بگیرید
5 + 10
به توکن اول که عدد 5 است میرسیم و در لیست متدهای یونری سرچ میکنیم. عدد 5 که یک توکن از نوع
INT
است پیدا میشود و مقدار متغیر لفت برابر با ۵ خواهد بود
حال به ادامه پردازش فانکشن اصلی برمیگردیم
توکن بعدی که به آن برخورد میکنیم توکن جمع + است
این توکن در بین متدهای باینری پیدا میشود
پس ما شروع به پردازش ادامه عبارت میکنیم یعنی عدد 10
و متغیر لفت را به عنوان ورودی به متد پردازش جمع داده میشود
درواقع برنامه ما به این شکل اجرا میشود
left = 5
parseBinaryOperation(5)
نحوه پردازش متد
parseExpression
بصورت بازگشتی است و ممکن است درک آن کمی سخت باشد ولی با چند بار تحلیل آن به نحوه کاربرد آن پی میبرید

نکته دیگری که لازم به ذکر است این است که در شروع برنامه اولین توکن توسط این متد چک میشود
و به محض پیدا شدن یک متد مرتبط با آن توکن کنترل به آن متد سپرده میشود و تعدادی توکن توسط این
متدها خوانده میشود و ممکن است توکن بعدی که توسط این متد خوانده شود چندین توکن بعد از توکن اول باشد
برای مثال برنامه زیر را در نظر بگیرید
در اینجا اولین بار و دومین باری که این متد  در متد قبلی یعنی متد
parse
فراخوانده میشود. نشان داده شده است.
این متد بارها و بارها بصورت بازگشتی توسط متدهای پارسی که آن را فراخوانی میکنیم صدا زده میشود.

1-> fn hello() {

	print("Hello World")

}
2-> hello()
*/
func (p *parser) parseExpression(priority int) Node {
	var left Node
	if fn, ok := p.unaryOperations[p.currentToken.Type]; ok {
		left = fn()
	}

	/*
		این بخش برای چک کردن این است که آیا توکن بعد از لفت توکن باینری است یا نه
		و اینکه اولویت توکنی که با آن متد اصلی یعنی
		parseExpression
		صدا زده شده است کمتر از توکن پیش رو است یا خیر
		اجازه بدهید با مثال قضیه را روشن کنیم
		عبارت زیر را در نظر بگیرید
		2 * 3 + 5

		(اجرای اول)
		متد
		parseExpression
		با کمترین اولویت اجرا میشود یعنی
		priority = 1
		به عدد 2 میرسیم و پردازش میکنیم و درون متغیر لفت قرار میدهیم
		left = 2
		توکن بعدی عملیات ضرب است. عملیات ضرب اولویت بیشتری از مقدار اولویت اولیه یعنی متغیر
		priority
		دارد پس وارد حلقه میشویم. متد متناسب با عملیات ضرب را پیدا میکنیم که متد
		parseBinaryOperation
		است و مقدار لفت یعنی 2 را به آن پاس میدهیم
		بعد از اجرای این متد سمت چپ معادله 2 است و سمت راست را با کال کردن
		parseExpression
		و پاس دادن اولویت عملگر فعلی که ضرب هست بدست می آوریم
		(اجرای دوم)
		این متد عدد بعدی یعنی 3 را به عنوان لفت پردازش میکند
		چک میکند که اولویت توکن بعدی یعنی ضرب از جمع بیشتر است؟ جواب خیر است پس وارد حلقه نمیشود
		و متغیر لفت که مقدار 3 را دارد برمیگرداند. در اینجا برمیگردیم به دومین اجرای فانکشن
		parseExpression
		مقدار لفت و رایت مشخص شده و برمیگردیم به اجرا اول متد
		parseExpression
		با برگشت به این متد مقدار متغیر لفت تغییر کرده و تبدیل شده است به
		(2*3)
		حال که مقدار لفت را داریم سراغ حلقه بعد از لفت در اولین اجرای متد میرویم
		در اولین اجرا مقدار
		priority = 1
		است. چک میکنیم که توکن بعدی یعنی جمع از متغیر اولویت بیشتر است؟ جواب بله است
		پس وارد حلقه میشویم و لفت را به
		parseBinaryOperation
		پاس میدهیم
		در واقع داریم:
		parseBinaryOperation( (2*3) )
		در ادامه مسیر مثل قبل سمت راست معادله را بدست می آوریم که عدد 5 است
		و در نهایت خواهیم داشت
		( (2*3) + 5 )

		اما اگر عبارت ما
		2+3*5
		باشد چه اتفاقی می افتد؟
		این بار خلاصه تر پیش میریم
		(اجرای اول)
		left = 2
		priority(+) > 1  // YES
		parseBinaryOperation( 2 )
		2 + parseExpression( priority(+) )

		lowest_priority = 1
		parseExpression (lowest_priority) {
			left = 2
			priority(+) > lowest_priority  // YES
			parseBinaryOperation(2) {
				p = priority(+)
				2 + parseExpression( p ) {
					left = 3
					nextTokenPriority = *
					* > + // YES
					parseBinaryOperation(3) {
						left = 3
						operator = *
						p = priority(*)
						right = parseExpression(p) {
							left = 5
							nextTokenPriority = lowest_priority
							lowest_priority > p // NO
							return left
						}

					}
				}
			}
		}
		نتیجه:
		( 2 + (3 * 5) )
		میدونم خیلی پیچیده شده الگوریتم های بازگشتی همینطور هستند ولی با چندبار مرور و گذاشتن
		بریک پوینت توی برنامه و تریس کد متوجه نحوه کارکرد الگوریتمش میشید


	*/
	pp := p.getPriority(p.nextToken.Type)
	for pp > priority {
		fn, ok := p.binaryOperations[p.nextToken.Type]
		if !ok {
			return left
		}
		p.next()
		left = fn(left)
	}
	return left
}

/********************************************************************************
* در زیر لیست تمام متد هایی که برای پارس کردن هر توکن نیاز است را مشاهده میکنید *
********************************************************************************/

/*
ساده ترین توکنی که میتوانیم پردازش کنیم توکن رشته است
ساختار استراکت رشته به این شکل است که یک متغیر به نام
Value
از نوع رشته در خودش دارد که مقدار آن توکن را در بر میگیرد
مثال
"Hello World!"
*/
type String struct {
	Value string
}

func (p *parser) parseString() Node {
	return String{Value: p.currentToken.Value}
}

/*
توکن بعدی که نسبتا تفسیر ساده ای دارد توکن آیدنتیفایر است
خود توکن را درون استراکت ذخیره میکنیم.
از آنجایی که آیدنتیفایر ها میتواند نام متغیر یا نام فانکشن باشند یک فلگ در نظر گرفته ایم
که متوجه شویم این توکن نام یک متغیر/ثابت است یا نام یک فانکشن
برای مثال
name
hello()
اولی نام متغیر
و دومی نام یک فانکشن است
*/
type Identifier struct {
	Token    Token
	Function bool
}

func (p *parser) parseIdentifier() Node {
	return Identifier{Token: *p.currentToken, Function: p.nextToken.Type == LPARENT}
}

/*
اگر به یک توکن از نوع عدد صحیح برخورد کردیم ابتدا باید آن را از رشته تبدیل به یک عدد کنیم
سپس مقدار عددی آن را درون استراکت مربوط به عدد صحیح قرار دهیم
*/
type Int struct {
	Token Token
	Value int
}

func (p *parser) parseInt() Node {
	v, _ := strconv.Atoi(p.currentToken.Value)
	return Int{
		Token: *p.currentToken,
		Value: v,
	}
}

/*
همان کار بالا را برای توکن های از نوع اعشار و بولین انجام میدهیم
*/
type Float struct {
	Token Token
	Value float64
}

func (p *parser) parseFloat() Node {
	v, _ := strconv.ParseFloat(p.currentToken.Value, 64)
	return Float{
		Token: *p.currentToken,
		Value: v,
	}
}

type Bool struct {
	Token Token
	Value bool
}

func (p *parser) parseBool() Node {
	b, _ := strconv.ParseBool(p.currentToken.Value)
	return Bool{
		Token: *p.currentToken,
		Value: b,
	}
}

/*
برای پارس کردن توکن
return
نیازی به خود توکن نداریم پس آن را ذخیره نمیکنیم
فقط قسمت جلوی کلمه کلیدی ریترن برای ما مهم است که
آن را درون یک متغیر درون استراکت ذخیره میکنیم
مثال
return 5
return a + b + c
return getScore()
از آنجایی که روبری کلمه کلیدی ریترن میتواند یک عبارت باشد پس ما متد
parseExpression
را فراخوانی میکنیم تا هرچیزی که جلوی این کلمه کلیدی است به عنوان مقدار در استراکت
مربوط به این کلمه کلیدی ذخیره شود
*/
type Return struct {
	Value Node
}

func (p *parser) parseReturn() Node {
	p.next() // از روی کلمه کلیدی ریترن میپریم
	return Return{
		Value: p.parseExpression(LOWEST_PRIORITY),
	}
}

/*
برای پارس کردن یک متغیر استراکتی تعریف میکنیم که حاوی دو متغیر است که هر دو از نوع
Node
هستند. مقدار
Name
یک
Identifier
است و
Value
هم یک عبارت
مثال
name         =     "Hello"
  |					  |
Identifier		  Expression



  c         =       a + b
  |					  |
Identifier		  Expression

این متد یکی از متدهای باینری است به محض اینکه به توکن مساوی میرسیم این متد فراخوانی میشود
و متغیر لفت به عنوان ورودی به این متد پاس داده میشود و درون متغیر نام قرار میگیرد
و قسمت جلوی مساوی به عنوان یک عبارت پردازش میشود
در واقع این متد به شکل زیر فراخوانی میشود
parseVariable(left)
		|
		|
{ Name: left, Value: parseExpression() }
*/

type Variable struct {
	Name  Node
	Value Node
}

func (p *parser) parseVariable(node Node) Node {
	p.next() // از روی توکن = میپریم
	ret := Variable{
		Name: node,
	}
	ret.Value = p.parseExpression(LOWEST_PRIORITY)
	return ret
}

/*
برای پارس کردن عباراتی که با پرانتز محصور شده اند نیازی به تعریف یک استراکت جدید نیست
کافی است مقدار درون پرانتز را به عنوان یک عبارت پردازش کنیم
برای مثال:
a = (b+c) * d
*/
func (p *parser) parseGrouped() Node {
	p.next() // از روی پرانتز باز میپریم
	exp := p.parseExpression(LOWEST_PRIORITY)
	p.next() // از روی پرانتز بسته میپریم
	return exp
}

/*
برای پردازش عملگرهای مانند منفی یا علامت تعجب که اگر قبل از یک عبارت بیایند معنی آن را تغییر میدهند بصورت زیر عمل میکنیم
خود توکن را درون استراکت ذخیره میکنیم و هرچیز که در مقابل توکن آمد را پردازش میکنیم و درون متغیر
Exp
قرار میدهیم
مثال:
-5
!true
*/

type UnaryOperator struct {
	Token Token
	Exp   Node
}

func (p *parser) parseUnaryOperator() Node {
	ret := UnaryOperator{
		Token: *p.currentToken, // ابتدا توکن را (- یا !) درون این متغیر ذخیره میکنیم
	}
	p.next()                            // حال که توکن را درون استراکت ذخیره کردیم از آن عبور میکنیم
	ret.Exp = p.parseExpression(PREFIX) // عبارت مقابل توکن را به عنوان یک عبارت پردازش میکنیم و درون این متغیر قرار میدهیم
	return ret
}

/*
برای پارس کردن عبارت های باینری استراکتی در نظر میگیرم که سمت چپ و راست یک عبارت است
و یک توکن که در بر گیرنده نوع عملگر است
مثال

	           ( a + b )                    +                 c
	               |                        |                 |
		   /	   |	   \       		 Operator           RIGHT
	      a        +         b
	      |        |         |
	     Left   Operator   Right

		          LEFT

همانطور که می‌بینید در عبارت
a + b + c
چون پردازش عبارت ها بصورت بازگشتی صورت میگیرد
عبارت
a + b
خود یک عبارت از نوع
BinaryOperator
است که درون متغیر
left
قرار میگیرد و متغیر
c
درون متغیر
right
قرار میگیرد
*/
type BinaryOperator struct {
	Token Token
	Left  Node
	Right Node
}

func (p *parser) parseBinaryOperation(left Node) Node {
	ret := BinaryOperator{
		Token: *p.currentToken,
		Left:  left,
	}
	/*
		برای آنکه بدانیم اولویت این عملگر در کل عبارت چه حق تقدمی دارد
		ابتدا اولویت توکن فعلی را محاسبه میکنیم سپس این اولویت را به متد
		parseExpression
		پاس میدهیم تا در آنجا در ساخت عبارت با اولویت صحیح درست عمل کند
	*/
	priority := p.getPriority(p.currentToken.Type)
	p.next() // از روی توکن عملگر میپریم
	ret.Right = p.parseExpression(priority)
	return ret
}

/*
پارس کردن یک بلاک کد بسیار ساده است هر وقت به یک آکولاد باز } رسیدیم
این متد را فراخوانی میکنیم و هر خط را به عنوان یک عبارت پردازش کرده و درون اسلایس
عبارت ها قرار میدهیم
*/
type Block struct {
	Statements []Node
}

func (p *parser) parseBlock() *Block {
	block := &Block{
		Statements: []Node{},
	}
	for p.nextToken.Type != RCURLY {
		p.next()
		exp := p.parseExpression(LOWEST_PRIORITY)
		block.Statements = append(block.Statements, exp)
	}
	p.next()
	return block
}

/*
استراکتی که برای پارس کردن دستورات شرطی در نظر گرفتیم دارای سه قسمت است
قسمت شرط‌، دستوارتی که بعد از صحت شرط اجرا میشوند و دستوراتی که درصورت عدم صحت شرط
اجرا میشوند. قسمت سوم اختیاری است و با کلمه کلیدی
else
مشخص میشود
*/
type If struct {
	Condition Node
	True      *Block
	Else      *Block
}

func (p *parser) parseIf() Node {
	p.next() // پریدین از روی توکن if
	exp := If{
		Condition: p.parseExpression(LOWEST_PRIORITY),
	}
	/*
	   در اینجا به انتهای قسمت شرط رسیده ایم و توکن بعدی کروشه باز است
	   از روی آن میپریم و بلاک کد را پردازش میکنیم
	*/
	p.next()
	exp.True = p.parseBlock() // بلاک کدی که در صورت صحت شرط اجرا میشود

	/*
		اگر توکن بعدی
		else
		نبود پس کار ما در اینجا به اتمام رسیده
	*/
	if !p.isNextToken(ELSE) {
		return exp
	}

	p.next()                  // از روی توکن الس میپریم
	exp.Else = p.parseBlock() // بلاک بعدی را پارس میکنیم
	return exp
}

/*
دستور حلقه از سه بخش تشکیل شده
شمارنده یا شمارنده ها

در زبان کهربا میتوان یک شمارنده یا دو شمارنده پاس داد
اگر یک شمارنده قرار دهیم مثال
for i in ["A","B","C"] { }
مقدار شمارنده برابر با هر کاراکتر خواهد بود اما اگر دو شمارنده به کار ببریم اولی مقدار اندیس
و دومی مقدار آرایه را در برمیگیرد

موضوع شمارش: که میتونه آرایه ، مپ یا یک رشته باشه

دستورات حلقه که یک بلاک کد هست
*/
type For struct {
	Key     *Identifier
	Value   *Identifier
	Subject Node
	Block   *Block
}

func (p *parser) parseFor() Node {
	p.next() // پریدن از روی توکن for
	key := p.parseIdentifier().(Identifier)
	ret := For{
		Key: &key, // پردازش متغیر شمارنده
	}

	/*
		چک کردن وجود شمارنده دوم
	*/
	if p.isNextToken(COMMA) {
		p.next() // از روی کاما میپریم
		value := p.parseIdentifier().(Identifier)
		ret.Value = &value
	}

	p.next()                                         // پریدن از روی متغیر شمارنده
	p.next()                                         // پریدن از روی توکن in
	ret.Subject = p.parseExpression(LOWEST_PRIORITY) // پردازش موضوع شمارنده
	p.next()                                         // پریدن از روی توکن بعدی و قرار گرفتن روی کروشه باز
	ret.Block = p.parseBlock()                       // پردازش بلاک کد
	return ret
}

/*
اگر به توکن ..  رسیدیم یعنی در حال پردازش یک عبارت رنج هستیم
استراکتی که برای پردازش این دستور در نظر گرفته شده شامل ابتدا و انتهای رنج هست
و قدمهای که اضافه میشه
قسمت قدم ها اختیاری است و قدم پیش فرض عدد ۱ هست
*/
type Range struct {
	From Node
	To   Node
	Step Node
}

func (p *parser) parseRange(from Node) Node {
	ret := Range{
		From: from, // توکن قبلی که یک عدد است را درون این متغیر قرار میدهیم
	}
	p.next()                                    // از روی ..  می‌پریم
	ret.To = p.parseExpression(LOWEST_PRIORITY) // عدد دوم را پردازش میکنیم و درون متغیر دوم قرار میدهیم

	/*
		در دستور رنج تعداد قدم های شمارش اختیاری است مانند
		1..10:2
		یعنی از 1 تا 10 ولی با قدم های دوتایی
		پس چک میکنیم که آیا توکن بعدی کولون : است یا نه
	*/
	if p.isNextToken(COLON) {
		p.next()                                      // از روی توکن : میپریم
		ret.Step = p.parseExpression(LOWEST_PRIORITY) // تعداد قدم های شمارش را در این متغیر قرار میدهیم
	}

	return ret
}

/*
برای پردازش دستور چاپ استراکت ما دو متغیر دارد که یکی نشان دهنده
آرگومان هاییست که به فانکشن پرینت ارسال میکنیم و دیگری نشان دهنده اینه که
در حال پردازش کدام یک از حالت های دستور چاپ هستم
print یا println
*/
type Print struct {
	Args    []Node
	NewLine bool
}

func (p *parser) parsePrint() Node {
	ret := Print{
		Args:    []Node{},
		NewLine: p.currentToken.Type == PRINTLN, // آیا توکن ما PRINTLN است ؟
	}
	p.next()                         // از روی توکن پرانتز باز میپریم
	ret.Args = p.parseFunctionArgs() // لیست آرگومان ها را پردازش میکنیم
	return ret
}

/*
این متد برای پردازش آرگومان ها استفاده می‌شود و موارد پردازش شده را بصورت
یک اسلایس از نود برمیگرداند
*/
func (p *parser) parseFunctionArgs() []Node {
	ret := []Node{}
	for p.nextToken.Type != RPARENT {
		p.next()
		if p.currentToken.Type == COMMA {
			p.next() // از روی توکن کاما میپریم
		}
		ret = append(ret, p.parseExpression(LOWEST_PRIORITY))
	}
	p.next()
	return ret
}

/*
آرایه ها در زبان کهربا میتوانند لیستی از همه چیز باشند
برای مثال
mylist = [ 1, "kahroba", 3.14, hello() ]
پس استراکتی که برای پردازش آرایه ها نیاز داریم فقط یک اسلایس از نود ها هست
*/
type Array struct {
	Nodes []Node
}

func (p *parser) parseArray() Node {
	p.next() // از روی براکت باز میپریم
	ret := Array{
		Nodes: make([]Node, 0), // اسلایس نود ها را مقدار دهی اولیه میکنیم
	}
	/*
		تا وقتی به انتهای آرایه نرسیدیم هر عبارت را پردازش میکنیم و داخل اسلایس قرار میدیم
	*/
	for p.currentToken.Type != RBRACKET {
		ret.Nodes = append(ret.Nodes, p.parseExpression(LOWEST_PRIORITY))
		p.next()
		if p.currentToken.Type == COMMA {
			p.next()
		}
	}
	return ret
}

/*
برای اینکه بتوانیم به اندیس یک آرایه یا مپ دسترسی داشته باشیم نیاز به تعریف
استراکت جدیدی داریم که در برگیرنده آرایه یا مپ مقصد و اندیس مورد دسترسی باشد
براث مثال:
a = [1,2,3,4]
print(a[0]) // 0
یا
a = {"name":"Kahroba","version":"0.01"}
print(a["name"],a["version"])
*/
type ArrayMapIndex struct {
	Subject Node // آرایه یا مپ مقصد را نگهداری میکند
	Index   Node
}

/*
از آنجایی که این یک متد باینری هست نام آرایه یا مپ قبلا درون متغیر لفت پردازش شده
و به این متد پاس داده میشه
*/
func (p *parser) parseArrayIndex(left Node) Node {
	p.next() // از روی براکت باز می‌پریم
	ret := ArrayMapIndex{
		Subject: left,
		Index:   p.parseExpression(LOWEST_PRIORITY),
	}
	p.next() // از روی براکت بسته میپریم
	return ret
}

/*
پردازش مپ هم مشابه پردازش آرایه است با این تفاوت که بجای اسلایس از یک مپ استفاده میکنیم
*/
type Map struct {
	Nodes map[Node]Node
}

func (p *parser) parseMap() Node {
	p.next()                           // از روی براکت باز می‌پریم
	ret := Map{Nodes: map[Node]Node{}} // استراکت را مقدار دهی اولیه میکنیم
	for {
		key := p.parseExpression(LOWEST_PRIORITY) // مقدار کی را پردازش می‌کنیم
		p.next()                                  // به توکن بعدی که : است میرویم
		p.next()                                  // از روی : میپریم
		val := p.parseExpression(LOWEST_PRIORITY) // مقدار ولیو را پردازش می‌کنیم
		ret.Nodes[key] = val                      // مقدار کی و ولیو را درون مپ استراکت قرار میدهیم
		/*
			اگر توکن بعدی کاما بود دو توکن میخوانیم تا از روی کاما رد شویم
		*/
		if p.nextToken.Type == COMMA {
			p.next()
			p.next()
		}
		/*
			اگر به براکت بسته رسیدیم از روی آن میپریم و کار پردازش مپ را تمام میکنیم
		*/
		if p.nextToken.Type == RCURLY {
			p.next()
			break
		}
	}
	return ret
}

/*
برای پردازش فانکشن از استراکت زیر استفاده میکنیم
نام بصورت رشته ای ذخیره میشود. لیست پارامترها که بصورت اسلایس از آیدینتیفایر است و بلاک اصلی کد
*/
type Function struct {
	Name   string
	Params []*Identifier
	Body   *Block
	Scope  *Scope
}

func (p *parser) parseFunction() Node {
	p.next() // از روی توکن fn می پریم
	exp := Function{Name: p.currentToken.Value}
	p.next() // از روی توکن نام فانکشن میپریم
	exp.Params = p.parseFunctionParams()
	exp.Body = p.parseBlock()
	return exp
}

/*
عملکرد این فانکشن هم مانند فانکشن
parseFunctionArgs
است و برای پارس کردن پارامتر های یک فانکشن از آن استفاده میکنیم
*/
func (p *parser) parseFunctionParams() []*Identifier {
	ret := []*Identifier{}

	for p.nextToken.Type != RPARENT {
		p.next()
		if p.currentToken.Type == COMMA {
			p.next()
		}
		ret = append(ret, &Identifier{Token: *p.currentToken})
	}
	p.next()
	p.next()

	return ret
}

/*
استراکت قبلی برای تعریف فانکشن بود
برای صدا زدن فانکشن هم نیاز به یک استراکت جدید داریم که بدانیم
فانکشن ما با چه آرگومانهایی صدا زده شده است.
*/
type FunctionCall struct {
	Function Node
	Args     []Node
}

func (p *parser) parseFunctionCall(function Node) Node {
	return FunctionCall{
		Function: function,
		Args:     p.parseFunctionArgs(),
	}
}

/*
این استراکت برای تغییر مقدار دو متغیر توسط فانکشن
swap
استفاده میشود
A و B
ورودی های فانکشن هستند
*/
type Swap struct {
	A Node
	B Node
}

func (p *parser) parseSwap() Node {
	p.next() // از روی توکن swap میپریم
	p.next() // از روی توکن ) میپریم
	ret := Swap{
		A: p.parseExpression(LOWEST_PRIORITY), // متغیر اول را پردازش میکنیم
	}
	p.next()                                   // از روی متغیر اول میپریم
	p.next()                                   // از روی کاما میپریم
	ret.B = p.parseExpression(LOWEST_PRIORITY) // متغیر دوم را پردازش میکنیم
	return ret
}

/*
این استراکت برای دریافت ورودی های
import
استفاده می شود
*/
type Import struct {
	Filename Node
}

func (p *parser) parseImport() Node {
	p.next() // از روی توکن import میپریم
	p.next() // از روی ) می پریم
	ret := Import{
		Filename: p.parseExpression(LOWEST_PRIORITY), // نام فایل را پردازش میکنیم
	}
	p.next()

	return ret
}

/*
این استراکت برای دریافت ورودی استفاده میشود
*/
type Input struct {
	Promp Node
}

func (p *parser) parseInput() Node {
	p.next() // از روی توکن input میپریم
	p.next() // از روی توکن ) میپریم
	ret := Input{
		Promp: p.parseExpression(LOWEST_PRIORITY), // promp را پردازش میکنیم
	}
	p.next()
	return ret
}

/*
این استراکت برای نگهداری طول یک آرایه یا مپ استفاده میشود
*/
type Len struct {
	ArrMap Node
}

func (p *parser) parseLen() Node {
	p.next()                                               // از روی توکن len میپریم
	p.next()                                               // از روی توکن ) میپریم
	ret := Len{ArrMap: p.parseExpression(LOWEST_PRIORITY)} // نام مپ یا آرایه را پردازش میکنیم
	p.next()                                               // از روی ( میپریم
	return ret
}
