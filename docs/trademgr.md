
## **trademgr**

Services
- *start thread for each symbol and strategy combination* 
- *thread exits only if trade is completed*
- *invoke timebased/continous signals to check trade confirmation from algotrading-analysis-service*
- *initiate order with zerodha*
- *monitor exit conditions*

Component
* [trademgr-signals](/trademgr-signals.html) 

### Overall

```plantuml
@startuml
participant TradeMgr as trademgr #green
participant tm.Signals as tm.signals #lightgreen
entity apiclient as apiclient #darkred
database    TimeScaleDB    as tsdb #grey


?-> trademgr ++ : ""->""\n**start** trademgr
trademgr -> tsdb : Fetch Enabled UserStrategies
tsdb --> trademgr : UserStrategies
|||
    loop per symbol
        trademgr -> tm.signals : start thread
        activate tm.signals #FFBBBB
        loop
            tm.signals -> apiclient : check signals
        end
    tm.signals -> tm.signals : await signal/trade
    tm.signals --> trademgr : Trade completed/Signal not found
    |||
    end
[-> trademgr --++ : ""->""\n**stop** trademgr
|||
trademgr -> tm.signals : check signal 
@enduml





