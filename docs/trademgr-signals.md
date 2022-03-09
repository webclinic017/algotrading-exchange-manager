
```plantuml
@startuml
scale 350 width
[*] --> scanSignals
state scanSignals {
[*] --> NewSymbol : Register Scanner
NewSymbol --> TimeTrigerred  
NewSymbol --> ContinousScan 
}
@enduml