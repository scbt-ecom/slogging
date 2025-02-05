# Introduction #
## Overview | Research ##
Пакет для распределенного логирования и трассировки с помощью встроенного пакета slog в Graylog с поддержкой всех используемых средств связи (AMQP, gRPC, HTTP). Под поддержкой имеется ввиду реализацию возможности передачи X-B3 заголовков с помощью встроенных методов со стороны сервера, а также реализацию чтения этих же заголовков со стороны клиента, и дальнейшей передаче их в контекст и создания экземпляра логера с дополнительными X-B3 полями.  

Также в процессе разработки было пересмотрено отношение к X-B3 заголовкам, и так как у нас нет полноценной системы централизованной трассировки, а эту роль выполняет Graylog было принято решение избавиться от всех X-B3 заголовков кроме X-B3-TraceId за ненадобностью и бесполезностью прочих. 

Изучение X-B3 трассировки по Zipkin дало понять, что до этого мы неправильно генерировали X-B3-TraceId, в оригинале это должно быть 64 или 128 битное шестнадцатеричное число, у нас же в некоторых местах использовался UUID (исправлено). (https://github.com/openzipkin/b3-propagation)

Был пересмотрен и подробно изучен вопрос передачи логгера для комфортной и быстрой работы с ним, решением стало передача логгера в контексте, выполняющем утилитарную роль.
В следующей статье подробно поднят этот вопрос, ознакомтесь если возникнут вопросы
(https://www.kaznacheev.me/posts/en/where-to-place-logger-in-golang/)  
Теперь необходима передача контекста на всем пути работы программы, несмотря на то, что у нас это не всегда использовалось, это правильный подход.

Самостоятельная реализация и актуализация отправки по прикладному протоколу GELF данных в Graylog позволила уменьшить размер payload, путем устранения deprecated GELF data, которая потеряла актуальность в связи с обновлением протокола GELF и технологии Graylog, соответственно повысилась производительность, помимо этого сам пакет для логирования slog выйгрывает перед логрусом в 10 раз по скорости. (https://github.com/betterstack-community/go-logging-benchmarks?tab=readme-ov-file)

Был реализован утилитарный функционал для пакета slog, позволяющий создать fanout handler, так как slog не умел это из коробки, так же присутствует возможность гибкой настройки логгера как в стандартный вывод, так и в Graylog. Также был затронут UDP Batching, честно говоря я хотел переписать его с нуля, но когда побольше разобрался в этой теме, понял что он и так неплохо реализован, для этого был использован вспомогательный пакет slog-graylog. (https://github.com/samber/slog-graylog)  

## Сопоставление уровней пакета и Graylog ##
| slog    | Graylog |
|---------|---------|
| Debug   | 7       |
| Info    | 6       |
| Warn    | 4       |
| Error   | 3       |

Был рассмотрен вопрос добавления уровня Fatal с os.Exit(1), от реализации пришлось отказаться, как я считаю всех представленных уровней достаточно, а остальные лишь излишнее усложнение, к тому же разработчики пакета slog сами отказались от добавления этого уровня, и предложили писать os.Exit(1) там где это нужно (потерял ссылку на github issue). К тому же для реализации кастомного уровня необходима еще одна инкапсуляция из-за чего код может стать слишком сложным и поддерживаемым. В данный момент реализация представлена без лишней инкапсуляции.

# Usage #

## Getting ##
```bash
go get github.com/scbt-ecom/slogging@v1.0.0
```

## Initialization ##
```bash
log := slogging.NewLogger(
    slogging.InGraylog({graylogURL}, {graylogLogLevel}, {container_name}),
    slogging.SetLevel({logLevel}),
    slogging.WithSource(true),
    slogging.SetDefault(true),
)
```
Важно: в примере представлены все возможные опции, если ничего не указывать выставятся стандартные
## Описание опций настройки ##
|                     | Установлена                                                                                                                                                                                                                        | Не установлена                                                                           |
|---------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------|
| InGraylog(params)   | Отправка логов в Graylog по протоколу UDP по указанному адресу,  учитывая выставленный уровень логов, все логи отправляемые в Graylog будут иметь указанное дополнительное поле "container_name"                                   | Отправка логов в Graylog отсутствует                                                     |
| SetLevel(params)    | Установка указанного уровня логов для записи в стандартный поток вывода                                                                                                                                                            | Будет установлен стандартный уровень "info" для записи в стандартный поток вывода        |
| WithSource(params)  | Добавление к логам для записи в стандартный поток вывода полей, отвечающих за источник лога (функция, файл, строка)                                                                                                                | Дополнительные поля, отвечающие за источник лога отсутствуют в стандартном потоке вывода |
| SetDefault(bool)    | В случае true, логгер будет установлен стандартным, то есть при вызове slog.Info например из любой точки программы, будет использован инициализированный логгер, если была включена отправка в Graylog, она тоже будет произведена | В случае false будет использован стандартный логгер slog.Default при slog.Info например  |

## Обычное использование по шагам ##
1) Инициализируем в main.go наш логгер со всеми надстройками
2) Если нужен TraceID в запросе создаем Middleware/Interceptor для HTTP/AMQP/GRPC, логируем из контекста, передаем контекст всегда дальше по всем слоям вплоть до конца пути данных
3) Если не нужен TraceID, либо ставим SetDefault(true) и логируем просто со slog.Info(), slog.Warn() etc. Либо также передаем инициализированный по всей программе, не используя Middleware/Interceptor

### HTTP Middleware ###
При использовании этой middleware, к каждому запросу автоматически будет ставится TraceID, логгер с новым контекстом будет лежать в req.Context(), так как передача логгера теперь в контексте
```
tracemw := slogging.HTTPTraceMiddleware(log)

http.HandleFunc("/", tracemw(helloWorld))

// mux example
rules := r.Path("/rules").Subrouter()
rules.Handle("/", ruleGetExampleHandler)

rules.Use(slogging.MuxHTTPTraceMiddleware(log))
```

### GRPC Middleware ###
```
Возможность создать GRPC Interceptor присутствует, как только использую сразу сюда добавлю
```

### AMQP Middleware ###
```
Возможность создать AMQP Middleware присутствует, как только использую сразу сюда добавлю
```

## Логирование из контекста с логгером ##
В примере приведены также экстра поля для логирования в Graylog
```
slogging.L(ctx).Info("example log message",
		slogging.ErrAttr(errors.New("example error message")),
		slogging.StringAttr("hello", "world"),
		slogging.IntAttr("bye", 12),
		slogging.AnyAttr("data", object),
		slogging.FloatAttr("bye", 14.88),
		slogging.TimeAttr("timestamp", time.Now()),
	)
```
## Создание логгера с новыми полями для передачи дальше ##
```
log := slogging.L(ctx).With(
		slogging.StringAttr("module", "keycloak"))

ctx = slogging.ContextWithLogger(ctx, log)
// Передаем контекст дальше, TraceID в нем сохранится, в логгере добавится дополнительное поле
```

## Создание контекста с логгером ##
```
ctx = slogging.ContextWithLogger(ctx, slog.Default())
// Передаем контекст дальше, вставится стандартный логгер без TraceID, зато с логированием в Graylog
```

## Создание контекста с логгером со случайным TraceID (полезно для тестов) ##
```
ctx := slogging.Context()
```
