# ProtoKey

****

## Расположение

**protokey**
``
./protokey/protokey
``

**protocli**
``
./protokey/protocli
``

## Структура проекта protokey

```mermaid
classDiagram
    class Command {
        <<interface>>
        +GetType() OperationType
    }
    
    class SetCommand {
        +Key string
        +Value int
        +GetType() OperationType
    }
    
    class GetCommand {
        +Key string
        +GetType() OperationType
    }
    
    class ListKeysCommand {
        +Prefix string
        +GetType() OperationType
    }
    
    class Response {
        <<interface>>
        +GetType() OperationType
    }
    
    class SetResponse {
        +Err error
        +GetType() OperationType
    }
    
    class GetResponse {
        +Value int
        +Err error
        +GetType() OperationType
    }
    
    class ListKeysResponse {
        +Keys []string
        +Err error
        +GetType() OperationType
    }
    
    class Service {
        <<interface>>
        +Set(key string, value int) error
        +Get(key string) (int, error)
        +Keys(prefix string) ([]string, error)
    }
    
    class ProtoKeyService {
        -commands chan Command
        -responses chan Response
        +Set(key string, value int) error
        +Get(key string) (int, error)
        +Keys(prefix string) ([]string, error)
        -validateValue(value int) error
    }
    
    class Store {
        <<interface>>
        +ProcessCommand(cmd Command) Response
        +Load(commands []Command) error
        +PopProcessedCommands() []Command
    }
    
    class StoreImpl {
        -data map[string]int
        -processedCommands []Command
        -mu sync.RWMutex
        -processedCommandsMu sync.Mutex
        +ProcessCommand(cmd Command) Response
        +Load(commands []Command) error
        +PopProcessedCommands() []Command
    }
    
    class SnapshotService {
        <<interface>>
        +Append(commands []Command) error
        +GetSnapshot() ([]Command, error)
    }
    
    class FileSnapshotService {
        -filePath string
        -mu sync.Mutex
        -file *os.File
        -writer *bufio.Writer
        +Append(commands []Command) error
        +GetSnapshot() ([]Command, error)
        +Close() error
    }
    
    class Engine {
        -store Store
        -snapshotService SnapshotService
        -commandChan chan Command
        -responseChan chan Response
        -stopChan chan struct
        -snapshotDelay time.Duration
        -mu sync.Mutex
        +Run()
        -init() error
        -flushProcessedCommands()
    }
    
    class Handler {
        -service Service
        +GetValue(w http.ResponseWriter, r *http.Request)
        +SetValue(w http.ResponseWriter, r *http.Request)
        +ListKeys(w http.ResponseWriter, r *http.Request)
    }
    
    Command <|-- SetCommand
    Command <|-- GetCommand
    Command <|-- ListKeysCommand
    
    Response <|-- SetResponse
    Response <|-- GetResponse
    Response <|-- ListKeysResponse
    
    Service <|.. ProtoKeyService
    Store <|.. StoreImpl
    SnapshotService <|.. FileSnapshotService
    
    ProtoKeyService --> Command
    ProtoKeyService --> Response
    
    Engine --> Store
    Engine --> SnapshotService
    Engine --> Command
    Engine --> Response
    
    Handler --> Service
    
    StoreImpl --> Command
    StoreImpl --> Response
    
    FileSnapshotService --> Command
```


## Сценарий обработки set

```mermaid
sequenceDiagram
    participant Client
    participant Handler as HTTP Handler
    participant Service as ProtoKeyService
    participant Engine
    participant Store as StoreImpl
    participant Snapshot as FileSnapshotService

    Client->>Handler: POST /set {key, value}
    activate Handler
    
    Handler->>Service: Set(key, value)
    activate Service
    
    Service->>Engine: Send SetCommand
    activate Engine
    Note right of Engine: Command через channel
    
    Engine->>Store: ProcessCommand(SetCommand)
    activate Store
    
    Store-->>Engine: SetResponse
    deactivate Store
    
    Engine-->>Service: SetResponse через channel
    deactivate Engine
    
    Service-->>Handler: error/nil
    deactivate Service
    
    Handler-->>Client: HTTP Response
    deactivate Handler
    
    loop Периодическое сохранение
        Engine->>Store: PopProcessedCommands()
        activate Store
        Store-->>Engine: []Command
        deactivate Store
        
        Engine->>Snapshot: Append(commands)
        activate Snapshot
        Snapshot-->>Engine: error/nil
        deactivate Snapshot
    end
```

## Пример работы protocli + protokey

**Получение help protocli**
```
MBP-Maksim:protocli maksimveselov$ ./main 
A command line client for interacting with the Protokey key-value store

Usage:
  protocli [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  get         Get a value by key
  help        Help about any command
  keys        Get keys matching prefix
  set         Set a key-value pair

Flags:
  -h, --help   help for protocli

Use "protocli [command] --help" for more information about a command.
```

**Невалидная команда protocli**
```
MBP-Maksim:protocli maksimveselov$ ./main sq
Error: unknown command "sq" for "protocli"

Did you mean this?
        set

Run 'protocli --help' for usage.
unknown command "sq" for "protocli"

Did you mean this?
        set

MBP-Maksim:protocli maksimveselov$ 
```

**Set через protocli**
```
MBP-Maksim:protocli maksimveselov$ ./main set key 10
OK
```

**Set через protocli с невалидным значением**
```
MBP-Maksim:protocli maksimveselov$ ./main set key 123123123123123123
Error (400): bad request
```

**Get через protocli**
```
MBP-Maksim:protocli maksimveselov$ ./main get key
10
```

**Keys через protocli**
```
MBP-Maksim:protocli maksimveselov$ ./main keys k
key1
key
```