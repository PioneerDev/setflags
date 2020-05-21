# Set Flags
  
An exemplar project for setting personal, verifiable flags as part of regular investing in yourself.

## Introduction

A functioning prototype at https://group-mixin.droneidentity.eu.

= @todo Translate into English = 
新型立志红包让群众监督你完成自己立下的宏愿。
定投者发起以一年为限的定投红包，许愿每天干成什么，最少要多少人验证, 一共发放多少金额: 其中10%奖励给证人，其余是押金在每天完成任务后归还自己。
当全部任务顺利完成时，已经发给证人的奖励将以SFC（立志币）返还发起人。
跟传统红包不一样，定投者（红包发起人）在报告任务完成并获得成功验证以后，可以赢回押金。 为获知任务完成情况，系统要求定投者和证人（红包领取人）配合完成以下规定动作：
* 每天系统提醒定投者报告当天任务的完成情况。
* 根据定投者的报告，验证者独立判断该任务是否已经完成。
* 系统只奖励结论符合共识的证人。
* 发起过立志红包的证人获得更多奖励。
* 只有经验证成功完成了任务，系统才退还定投者当天押金。
* 每个红包至少价值1 USDT或者1 SFC。
* [optional] 消息根据用户语言自动翻译!
* [optional] 如果立志人不折不扣地完成全部任务，全部红包价值完璧归赵, 并免费获得全年完整的证据链 ! 

## Requirements

### Use Cases

## Design

We follow the model-view-controller architecture pattern to design the chatbot.

* Model

![Entity Relationship diagram](https://github.com/set-flags/set-flags/blob/master/docs/models.png)

NOTE. The diagram is created using PlantUML, [after the file is opened inside GitPod](https://gitpod.io/#https://github.com/set-flags/set-flags/blob/master/docs/models.puml), type "Alt + D" to edit the diagram.

* View

The following RESFTful API table is generated from [our Swagger specification](https://github.com/set-flags/set-flags/blob/master/docs/models.yml) using [the Swagger Editor](https://editor.swagger.io) by the [markdown-swagger tool](https://github.com/rmariuzzo/markdown-swagger).
<!-- markdown-swagger -->
 Endpoint                                     | Method | Auth? | Description                                                                                          
 -------------------------------------------- | ------ | ----- | -----------------------------------------------------------------------------------------------------
 `/flags`                                     | GET    | No    | list all the flags                                                                                   
 `/flags/{id}`                                | POST   | No    | create a flag                                                                                        
 `/flags/{id}/{op}`                           | PUT    | No    | update an existing flag with operations for verification (yes, no) after uploaded the evidence (done)
 `/myflags/{id}`                              | GET    | No    | list all flags of the user                                                                           
 `/flags/{flag_id}/evidences/{attachment_id}` | POST   | No    | upload evidence                                                                                      
 `/flags/{flag_id}/evidences`                 | GET    | No    | list all the evidences since yesterday                                                               
 `/users/{user_id}/rewards/{flag_id}`         | GET    | No    | check the total rewards received by the user for the flag                                            
 `/assets/{id}`                               | GET    | No    | get the asset information                                                                            
<!-- /markdown-swagger -->

* Controller

Front-end

## Implementation

* The source will be organised using Model-View-Controller (MVC) architecture.

models/

views/

controllers/
