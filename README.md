
一、环境依赖
	
	1.golang 
	2.mysql

二、使用说明

	1.填写config.yml
	2.修改 main.go 中 mysql 数据库帐户密码
	3.设置数据库名称 group_helper


三、命令

	命令一：领糖果

	命令二：创建社群#社群名称#总量#份数   （ 如：创建社群#吹牛逼社群#10000#100 ）

	命令三：公告#大家好

～～～～～～～～～～～～～～～～～～

注：总量和份数为数字，暂时只支持CNB

四、项目简介

基于mixin network实现一个社群助手

五、传统建立社群的方案

在电报群建立群，然后邀请用户加入，并发放糖果吸引用户入群

	 1.项目方发放糖果需要手续费（手续费问题）
	 2.技术成本等因素导致糖果不能及时发放（技术问题或者说团队时间成本问题）
	 3.用户加入群后不知道糖果什么时候发放没有邀请他人的欲望（扩大社群规模）
	 4.糖果的不及时发放会导致社群用户流失（信任问题）
	 5.电报注册很简单只需要手机号验证即可，所以存在大量羊毛团队

六、为什么mixin上适合做社群助手
	
	1.基于mixin network + messager 这种快捷、免费环境下能有效解决以上弊端
	2.并且通过mixin messager的注册机制（谷歌验证码），能有效减少羊毛党
	3.机器人能24验证用户是否完成项目方（通过api）的任务，并及时发放糖果

七、项目介绍
	
	项目方
	1.通过机器人建立社群  
	2.通过机器人支付糖果  
	3.通过机器人布置任务给用户（后期规划）
	4.给社群用户发布公告 
	5.持仓验证（后期规划）
	6.多级邀请机制（后期规划）
	7.支持mixin全币种建立社群（后期规划）
	

	用户
	1.加入社群就能找社群助手领取糖果
	2.每邀请一次，或者完成一次任务都能领取当前的糖果（后期规划）


项目方能给mixin messager带来更多的用户，促进mixin生态的发展

八、使用实例

创建社群
![Image text](https://raw.githubusercontent.com/ewnk/grouphelper/master/img/2.jpg)

创建后机器人会邀请你进群，并给予管理员权限
![Image text](https://raw.githubusercontent.com/ewnk/grouphelper/master/img/3.jpg)

邀请用户加入后用户可回复 “领糖果”
![Image text](https://raw.githubusercontent.com/ewnk/grouphelper/master/img/4.jpg)

管理员可以给自己的社群用户发送公告
![Image text](https://raw.githubusercontent.com/ewnk/grouphelper/master/img/5.jpg)

九、更新计划

	代码结构优化
	采用分布式数据库
	项目方定义任务
