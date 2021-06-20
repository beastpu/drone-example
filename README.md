# drone CI vs gitlab CI
# drone 
drone 是一个非常轻量级的持续集成工具，以yml的形式定义pipline，相比于jenkins pipline复杂的grovy语法，drone pipline非常简单易懂，学习成本低。  
虽然gitlab pipline也是yml定义，但是相比drone需要定义的参数更多，步骤繁琐。drone正是基于这种开发人员可以快速定制复杂的pipline来实现高效的CICD在繁忙的开发团队。  

## 局限性
drone是个轻量级工具，相对gitlab,jenkins，teamcity等工具提供的功能就不是很多。  
drone开发团队被收购以后推出了企业版和开源社区版。

社区版和企业版区别参考：<https://www.drone.io/enterprise/opensource/#features>

## example
.drone.yml inclues four piplines for different purposes.
 - first pipline build & test golang service when push event is triggered 
 - second pipline push package to remote repository.
 - last pipline deploy services fot differnet environment which are triggered by promotion event in manaul 
 
### 底层原理
```
steps:
- name: backend
  image: golang
  commands:
  - go build
  - go test
```
执行的脚本会转换为entrypoint。
set -e保证了根据状态码进行。  
set -x 进行变量的输出。
```
#entrypoint.sh
#!/bin/sh
set -e
set -x

go build
go test
```
最终转换的命令如下：
```
docker run --entrypoint=backend.sh golang
```

volume共享

执行pipline前首先会拉取git仓库，仓库会在pipline多个步骤间共享，即使步骤间使用不同的镜像。  
git clone时docker会创建一个volume存储代码。然后volume就会挂载到不同的容器。不会随着容器删除而删掉。只有pipline结束时才会被系统删除。

# gitlab CI
gitlab 虽然pipline相对复杂，但提供的功能丰富，比如可以通过`DAG`自由组合job间的依赖关系。使流程尽可能快的执行。
也可以通过条件规则（if only except等）来判断job是否执行，when:delayed实现延迟执行。不得不说功能非常强大。
gitlab页面也提供了可视化pipline的功能，即所见即所得。

## pipline
### basic piplines
.gitlab-ci.yml例子中是最基本的pipline,定义了三个阶段build->test->deploy.
相同阶段的任务会并行执行，如test-job1和test-job2.

### Directed Acyclic Graph
DAG能被用于job间的构建依赖，尽可能快的执行jobs.  
needs参数可以实现job间的依赖，使jobs非有序的执行，是否执行只与依赖的job有关。 

我们修改了`deploy-prod`，让它只依赖于`test job 1:3`。
也就是说当`test job 1:3`执行完deploy-prod会立即执行，而不需要等待`test job 2:3`. 

```
deploy-prod:
  stage: deploy
  needs: 
  - test job 1:3
  script:
    - echo "This job deploys something from the $CI_COMMIT_BRANCH branch."
```
### multi piplines
gitlab 支持多个pipline在一个项目中  
触发一个child piline
```
microservice_a:
  trigger:
    include: subProject/.gitlab-ci.yml
```
### group job
同一个stage的多个job会使可视化界面复杂，gitlab支持名字类似的job自动组成一个group。格式如下：
- A slash (/), for example, test 1/3, test 2/3, test 3/3.
- A colon (:), for example, test 1:3, test 2:3, test 3:3.
- A space, for example test 0 3, test 1 3, test 2 3.

```
test job 1:3:
  stage: test
  image: alpine:latest
  script:
    - echo "This job tests something"


test job 2:3:
  stage: test
  script:
    - echo "Thifs job tests something, but takes more time than test-job1."
    - echo "After the echo commands complete, it runs the sleep command for 5 seconds"
    - echo "which simulates a test that runs 20 seconds longer than test-job1"
    - sleep 5

test job 3:3:
  stage: test
  script:
    - echo "Thifs job tests something, but takes more time than test-job2."
    - echo "After the echo commands complete, it runs the sleep command for 5 seconds"
    - echo "which simulates a test that runs 20 seconds longer than test-job2"
    - sleep 5
```
### debug功能
当job没有按照期待的执行，gitlab 支持debug相应的job.  
当job运行时，点击右上角的debug按钮，可以进入到容器。  
If you have the terminal open and the job has finished with its tasks, the terminal blocks the job from finishing for the duration configured in [session_server].session_timeout until you close the terminal window.

### cache机制
git repo上创建的文件并不能在下一个job上共享，因此gitlab提供cache机制缓存文件在多个job中。
```
cache:
    paths: 
     - lib/
```
### job artifacts
job 可以输出一个归档的文件盒目录，这个输出叫做工件。
你可以下载工件通过gitlab ui或者api.  
pdf job执行命令输出mycv.pdf，一周内可以任意下载。
```
pdf:
  script: xelatex mycv.tex
  artifacts:
    paths:
      - mycv.pdf
    expire_in: 1 week
```
### 权限
最后说一下权限的问题，gitlab中repostory可以设置protected branch。  
假设master分支被保护，其他用户merge,push时都需要代码拥有者的支持。  
pipline运行在保护的分支，可以做到一定的权限管理。只有相应权限的人才可以触发pipline.  
> <https://docs.gitlab.com/ee/user/project/protected_branches.html> 

# 结论
如果产品的业务逻辑复杂，推荐使用gitlab,支持功能非常多，如果你发现某一功能不支持，一定是没有好好看文档。  
drone轻量，更适合二次开发。  

整理不易，求个👍 
