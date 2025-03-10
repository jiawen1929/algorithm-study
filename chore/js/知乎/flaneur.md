# 知乎经典回答摘录

https://www.zhihu.com/people/flaneur2023/answers?page=3

1. 为啥 redis 使用跳表(skiplist)而不是使用 red-black？
   写工程代码时候一般都是这个回路：

   - 在所有的选项里面找一个写起来最省事的
   - 跑通了
   - 跑着意外的还行
   - 没人抛问题
   - 不改了

2. 什么情况下需要考虑内存屏障？
   简单说：如果你玩多线程老老实实对共享资源加锁，甚至如果你不玩多线程，那你完全可以忘掉内存屏障这回事。它和你无关。
3. 面试时怎么样定义一个应届生的分布式系统设计能力？
   应届生就做做题算了，你们公司有几个正儿八经设计分布式系统的？2024 年你准备设计 k8s 还是 spark pytorch？
   现在是 ai 时代了，面试时应该更重视一个应届生的万卡集群大模型基础训练能力。
4. 为什么Rust的包管理器Cargo这么好用？
   之前看 dhh 的那本 getting real 里有句印象很深是说：成功才是成功之母。
   `cargo 早期的重要作者是 katz，他在 ruby 的包管理 bundler 已经成功了一次，到 rust 这边再成功一次理所当然。`
5. 为什么要使用环形缓冲区？
   ring buffer 基本上是经过几十年探索后高性能 ipc 的最优解了，这东西又简单又对 cache 友好，还对同步友好。
6. HTTPS能否完全取代HTTP？
   2024 年公网上不该有 HTTP 存在了。

   （根据 zero trust 的理论在内网也不该有）

7. 腾讯会议免费时长将缩短至 40 分钟，这对该公司未来发展有何影响？
   应该是一个收费功能，在有充分的文档同步上下文之后，会议确实不该超过 40 分钟。

   会议应当看作是 2PC 里的 Commit 阶段，在 Prepare 阶段做好充足的准备时，Commit 阶段应当快速完成。

8. 认真学习完 MIT 6.824和CMU 15445 可以找到相关的工作吗?
   个人偏见是做题路径过于清晰的领域的工作一定是更难找的。
9. CPU架构后续发展的方向是？
   感觉未来的 cpu 首先做好网卡和显卡的胶水层，跑好一个 io 密集负载就可以了。涉及到跑 ai，计算密集的含义已经变了，现在计算密集已经不算 cpu 能干的活了。这是计算能力的需求驱动的。再一个就是像 apple silicon 这样把一些加速模块嵌到 SoC 里，也是减少 cpu 本 u 的日常工作内容，是降功耗低碳驱动的。还有一个就是科技冷战中合规驱动的 riscv 这种开源信创概念，主打一个不想用也得用，不过还看不大清 riscv 运营的怎样，听说目前整的有点分裂。
10. rust 解决了什么问题？
    显然是包管理 ="=
11. 为什么工程中都用红黑树，而不是其他平衡二叉树？
    最近几年工程上都倾向用 BTree 做有序 map 了，内存局部性好。
12. 为什么 PostgreSQL 要使用多进程并发模型而不是多线程？
    单纯就是技术债太厚了没人敢碰。
13. 计算机上“中断”的本质是什么？
    就是个回调。
14. go 的 gc 性能如何？
    `go 的 gc 只为延时优化，没有分代，吞吐不如 jvm 系，体现就是没人拿 go 跑 big data 的负载。胜在简单灵活，光看延时一个指标，跑 rpc 负载无敌。`主要是 satb 的 barrier STW 短，不需要 remark。jvm 系期待 zgc 早日统一天下，别再背调参八股了。
15. 如何解决 async/await 的传染性？
    本来想说这个传染性是无解的，不过想了一下 etcd 中的 raft 状态机是一个这方面的典范，`它的出发点不是隔离 async/await 的传染性（go 里没有 async/await 传染性的问题），而是隔离所有副作用（网络、文件 IO），把无副作用的部分尽可能封在里面，把有副作用的地方都放在外面。同比到 raft-rs，你就可以看到这个包是没有 async 的，但是没有耽误你基于它配合 async 做一个 raft server。`这样抽出来一个**完全确定性且很厚的逻辑内核，把对外界交互的地方用很薄的一层对接上**。非常满足《A Philosophy of Software Design》中提到的 Deep Interface 的思想。只要 async 交互的界面足够小而且和主体逻辑独立无副作用，就不怕 async 的传染性。（想起来另一个例子是 clickhouse 的 processor 框架，计算密集的逻辑本体不会与 IO 相互侵入，IO 的内容都是调度器给它喂（push）的）
16. 有没有一句话可以媲美张载的「为天地立心，为生民立命，为往圣继绝学，为万世开太平」？
    文科世界的所有名言警句加起来也比不上瓦特一台蒸汽机。
17. 为什么说程序员不断的提高自己的技术有可能是一种误区？
    所谓技术深度和技术水平不能自动解决问题，反过来反而有很大概率是过拟合。
    解决真实世界需求的技术才是好技术。
    历史上的传奇程序员现在看是他们技术超神，可是当时他们也未必有知乎上的各位高手懂得多，但是他们动手解决了时代的需求，然后很残酷的现实是很多基础问题只需要解决一次，之后就不再需要第二个同类系统。
    Linus 当年就是操作系统领域有技术深度的专家吗，知乎上对 OS Kernel 各种犄角旮旯如数家珍的技术高手不比当年的 Linus 懂得多？知乎上对 redis 各种犄角旮旯如数家珍的技术高手懂的不比 Salvatore Sanfilippo 多？Salvatore Sanfilippo 来国内考他几道 redis 八股都不一定过的了二面。
    **程序员的价值只能通过解决全社会或者公司内部的新问题领域来体现，这里面往往有大量试错成分**。`在 well defined 和 already solved 的领域持续扎根考据并认为技术水平高反而可能是有害的，单论考据 GPT4 天下第一。`
18. 为什么程序员必须坚持写技术博客？
    我们做文本生成任务的，写博客就相当于做 finetuning。

    知识和能力不是通过 finetuning 获得的，但是 finetuning 可以把这部分知识和能力的表达进行激活，这样能够更好地对齐这些任务乃至泛化的任务。

19. 感觉资深程序员的优势不一定只在于代码输出能力强，而是在于`吸收大量优质代码语料和 Design docs 语料之后突变出来的推理能力强而且泛化能力好所以语料质量的区别就决定了程序员能力的区别`。（听说 ChatGPT 就是学了代码之后才学会的把复杂问题拆分成 N 个步骤分步解决不知道真假）
20. **linux内核复杂还是浏览器内核复杂？**
    投浏览器一票。
    内核的核心模块其实已经极简了，代码量高主要还是驱动堆出来的。`内核的核心模块是抽象的极致，本身是非常非常低熵的。但这里低熵不等于容易理解，难理解在于怎么去 get 到它的抽象。`
    但是浏览器是`需求的堆叠，各种 Web Standard 和野生标准的熵非常高`。光是理解它的行为就相当不容易了，反观内核的行为反而容易预期一些。
    可以做一个思想实验：做一个行为看起来像是 Linux 的内核（比如能跑个 Nginx 啥的给我两年应该可以苟出来），大概率会比做一个行为看起来像是 Chrome 的浏览器容易。
21. 为编程语言设计怎样的错误处理方式才是“好的”？
    Rust 这种 Result<R,E> 就挺好的，好像学名叫 Sum Types。感觉上接近错误处理的标准答案了。
    Result<R, E> 作为一个和类型，**强制开发者在编译时处理所有可能的情况（即成功和失败）。**这确保了错误不会被忽视，提高了代码的健壮性和可靠性。
22. 字节序问题
    字节序（Endianness）是指在计算机系统中，`多字节数据`（如整数、浮点数等）在内存中存储的字节顺序。
    utf-8不存在字节序问题
    UTF-8 是一种可变长度的字符编码，每个字符由1到4个字节组成。UTF-8 的设计使得每个字节在不同的字节序系统中都有明确的含义，无需依赖字节的顺序来解释数据。这种特性使得 UTF-8 在任何系统中都能一致地被解析和使用，无论系统是大端序还是小端序。
23. Spring，Django，Rails，Express这些框架技术的出现都是为了解决什么问题，现在这些框架都应用在哪些方面？
    **框架是在人们对组件技术失望之后出现的**，"do one thing and do it best" 和 "封装隐藏复杂度" 在项目里往往并不能如意，解开的耦合最终还是要集成，而**复杂度更多隐藏在组件与组件的集成，不在单个组件本身**，像讲复杂系统的书里经常提到的一句话是 "整体大于部分之和"，给我们小汽车的所有组件但我们组装不出小汽车。**框架做的事情，就是自带胶水，做一些脏活累活替我们集成这些组件，降低偶发复杂度，我们可以专心往里面填业务。**
24. 为什么干不掉 GIL
    GIL（Global Interpreter Lock，全局解释器锁） 是 Python 解释器（尤其是 CPython）中的一个关键机制，用于管理对解释器内部数据结构的并发访问。尽管 GIL 简化了多线程编程模型，但它也成为了 Python 在多核处理器上实现真正并行执行的一个瓶颈。
    - GIL（全局解释器锁） 是一个互斥锁（mutex），用于保护 CPython 解释器的内部数据结构，确保同一时间只有一个线程在执行 Python 字节码。这意味着，即使在多线程程序中，多个线程也无法真正并行地执行 Python 代码，而是必须依次执行。
    - 为什么存在 GIL？
      - 内存管理的简化：CPython 使用`引用计数作为主要的垃圾回收机`制，每个对象都有一个引用计数器，记录有多少引用指向它。当引用计数归零时，对象被销毁。为了确保引用计数的准确性，必须防止多个线程同时修改它。这在没有 GIL 的情况下，需要引入大量的锁机制来保护每个引用计数的修改，增加了复杂性和性能开销。
      - 线程安全：没有 GIL，需要在解释器的各个部分引入细粒度的锁，增加了出错的可能性，并且可能影响性能
      - 现有生态系统的兼容性
        大量的 C 扩展（如 NumPy、Pandas 等）依赖于 CPython 的内部实现，包括 GIL 的存在。移除 GIL 将需要重新设计这些扩展，导致兼容性问题，增加开发和维护成本。
25. 什么是软件事务内存（Software Transactional Memory, STM）
    - 软件事务内存（STM） 是一种用于`并发编程的抽象机制`，旨在简化多线程程序中的共享内存访问与同步。STM 通过借鉴数据库事务的概念，提供了一种更直观、安全且易于使用的方式来管理并发操作，避免传统锁机制中常见的问题，如死锁、竞态条件和锁粒度过细等。
    - STM 的核心思想是将一组内存操作封装在一个**事务（Transaction）**中，这些事务要么全部成功执行，要么全部回滚，保持数据的一致性。
    - STM 的实现方式
      1. 读集和写集：事务在执行过程中读取/写入的内存区域
      2. 版本控制：每个受保护的内存位置都维护一个版本号。事务在读取时记录当前的版本号，提交时检查版本号是否发生变化，以检测冲突。
      3. 乐观并发控制（Optimistic Concurrency Control）：STM 假设事务之间不会频繁冲突，因此允许事务并行执行，只有在提交时才进行冲突检测。
26. 拜占庭错误是不可克服的吗？为什么一些论文（如raft）都强调在非拜占廷错误情况下，都可以满足安全性，？

    - 内网的系统如果出现拜占庭问题一般会视为实现层面的 bug，在协议层面攻克它的成本太高了，这时候赶紧上备份修bug才是正经事。

    - 至于 bitcoin 这种 p2p 软件而言，对付拜占庭则属于系统生存下去的基本，野外随便哪个节点都有可能作弊。

27. 在开发过程中使用 git rebase 还是 git merge，优缺点分别是什么？
    团队里有新人的话建议统一 git merge，不然容易出来麻花，过去解完麻花再传授一番 git rebase 的原理然后 git push -f 一发又有什么收益呢，想回滚就找 build，想回溯历史的话 jira/wiki/gitlab 上有的是，个人感觉不是开源社区之类比较看重 git commit 来查 bug 和各种 backport 的话，团队项目这方面可以选择更不容易出错的用法
28. 云计算正在杀死运维吗?
    不管前端运维还是QA，不工程化都只有死路一条。
29. async/await异步模型是否优于stackful coroutine模型？
30. 当前主流的RPC框架有哪些？Dubbo与Apache thrift哪个更主流更容易上手？
    现在 rpc 框架只考虑 grpc 就可以了。
