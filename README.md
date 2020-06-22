# Set Flags
  
An exemplar project for setting personal goals to invest in yourself as regularly verifiable flags.

## Introduction

A functioning prototype at https://group-mixin.droneidentity.eu.

Everyone can customise the personal goal as a series of daily tasks by setting a flag. The purpose of issuing such a “flag” is to audit the achievements of the tasks which lead to the achievement of the goal eventually. For example, to publish a high-quality research paper, the goal can be decomposed into much smaller and more manageable tasks which are achievable every day, facilitating collegial friends to witness, applaud or verify the achievement. One may also use the concept for other personal goals such as physical exercises and training.

Specifically, the description of a flag includes the numbers of days till the deadline (365 days by default for New Years resolution), details of the task, minimal number of witnesses, and the total amount of prize, etc. Majority of the prize (for example 90%) will be rewarded back to the task taker when the task is successfully achieved by the end of each day. The remaining amount will be awarded to witnesses who "mine" the blockchain by verifying the tasks. Once the goal is fulfilled with all daily tasks achieved, the task taker can get all the prize returned, including those given away to the witnesses in the Set Flag Coins (SFC). Unlike traditional one-off prizes, which will not be rewarded back to the task taker, the flag can be returned fully by the end. 

## Requirements

To be aware of members’ progress, the system requires the task takers and the witnesses to do the following:
1. Download and install Mixin Messenger app, which is available in most app stores including Google Play Store and Apple App Store; 
2. Set up the Mixin Networks wallet using the telephone number and PIN as necessary; 
3. Search for the Mixin chatbot ID 7000103152;
4. You can “Set Flag” using the SFC, or any other cryptocurrency to create a goal to give away red packets in the due course. Ever since, each day the system will remind you of the todo tasks and ask for uploading evidences. If you find the reminders noisy or disturbing your flow, you can also mute it;
5. Depending on the evidence reported back to the system such as images, videos, etc., witnesses (i.e., other members in the group) will judge independently whether you have achieved the task on the day;
6. Such decisions are independent to each other, meaning that you have not been informed of other people’s decisions, hence a consensus mechanism is required to decide collectively on the task fulfilment;
7. The system will only reward the witnesses whose decision matches with the consensus. At the moment, such a consensus mechanism is fairly simple, i.e., the majority’s report rather than the minority report;
8. To encourage participation, the task takers who verify other people’s tasks can get slightly more reward than the witnesses who have never carried out any tasks;
9. The system rewards the daily amount only when the task taker has been verified to achieve the task;
10. The system rewards all the amount given to the witnesses as a final prize in SFC, only when the task takers can achieve the goal completely. The chain of custody, i.e. forensic evidence, for achieving the goal will be ready for download by the end.

## Design

We follow the model-view-controller architecture pattern to design the chatbot.

* Model

![Entity Relationship diagram](https://github.com/PioneerDev/setflags/blob/master/docs/models.png)

NOTE. The diagram is created using PlantUML, [after the file is opened inside GitPod](https://gitpod.io/#https://github.com/PioneerDev/setflags/blob/master/docs/models.puml), type "Alt + D" to edit the diagram.

* View

The following [RESFTful API](https://github.com/PioneerDev/setflags/blob/feature/rest-api/API-README.md) table is generated from [our Swagger specification](https://github.com/PioneerDev/setflags/blob/master/docs/models.yml) using [the Swagger Editor](https://editor.swagger.io) by the [markdown-swagger tool](https://github.com/rmariuzzo/markdown-swagger).
<!-- markdown-swagger -->
 Endpoint                             | Method | Auth? | Description                                                                                          
 ------------------------------------ | ------ | ----- | -----------------------------------------------------------------------------------------------------
 `/flags`                             | GET    | No    | list all the flags                                                                                   
 `/flag`                              | POST   | No    | create a flag                                                                                        
 `/flags/{id}/{op}`                   | PUT    | No    | update an existing flag with operations for verification (yes, no) after uploaded the evidence (done)
 `/myflags`                           | GET    | No    | list all flags of the user                                                                           
 `/attachments/{attachment_id}`       | POST   | No    | upload evidence                                                                                      
 `/flags/{id}/witnesses`              | GET    | No    | list all the witnesses                                                                               
 `/flags/{id}/evidences`              | GET    | No    | list all the evidences since yesterday                                                               
 `/me`                                | GET    | No    | Get the current user's information                                                                   
 `/auth`                              | POST   | No    | Upload the oauth token                                                                               
 `/users/{user_id}/rewards/{flag_id}` | GET    | No    | check the total rewards received by the user for the flag                                            
 `/assets/{id}`                       | GET    | No    | get the asset information                                                                            
<!-- /markdown-swagger -->

* Controller

Front-end

