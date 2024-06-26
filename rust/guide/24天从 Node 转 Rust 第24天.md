- 第 11 天：模块系统
  Rust 模块系统跟文件名、文件路径毫无关系
  模块就像是一个命名空间，用来把相似的事物放到一个篮子的。模块是可见性(visibility)的边界(barrier)，如 public, private 等。可以在一个文件内写多个模块。甚至可以在一个模块中嵌套模块。

  可以使用一些修饰语定制 pub 可见性，如下：

  pub(crate) - 在 crate 内部公开(public)，在外部则不可见。
  pub(super) - 仅对对父级模块公开。

- 第 12 天：字符串，第 2 部分

  函数参数应该用 &str 还是 String？
  第一个要问的问题是：“我应该拥有(own)还是借用(borrow)这个传入参数值？”

  对于这些还在困惑拥有(own)和借用(borrow)的人，可以这么想，“我需要拥有一个自有版本这个值，还是仅仅需要看一下它的数据而已？”

  - 当借用(borrow)参数值时 -> &str
  - 当需要一个自有的(owned)值时 -> String
  - 对于接受可能生成的字符串的函数 -> impl ToString

- 第 13 天：结果(Result) 和选项(Option)
  Rust 跟 JavaScript 不同，没有 undefined 和 null 的概念，因此 Rust 需要安全的表示“无”(nothingness)。这便是 Option 的使用场景。

  unwrap() 是一个简便的方法，但是如果 Option 是 None，它会 panic。
  unwrap_or() 可以提供一个默认值。
  unwrap_or_else() 可以提供一个闭包来生成默认值。
  unwrap_or_default() 返回 Default Trait 的默认值。

  在接收参数、返回参数时，相比于使用字符串、数字或者布尔值这种魔术值，**使用枚举会更好**，因为枚举可以表达更多的含义，不仅仅是 true 和 false 而已。

- 第 14 天：管理错误

  Rust 中的错误处理

- 第 15 天：闭包

  Fn, FnMut, 和 FnOnce 特质

  - Fn： 一个不可变的(immutably)借用其包进来的变量的函数
  - FnMut：一个可变的(mutably)借用其包进来的变量的函数
  - FnOnce: 一个消费其值（失去值的所有权），因此只能运行一次的函数，如：

  在 struct 内保存闭包

- 第 16 天：生命周期，引用，和 'static

  生命周期(lifetime, 或译生存期) 是 Rust 借用检查器内部的一个构造(construct)。每一个值，都有其创造之时和销毁之时。这就是它的生命周期。

  生命周期注解(lifetime annotation)是 Rust 语法，可以添加在一个引用之前，给该引用的生命周期一个命名标记。`当有多个引用时，且 Rust 无法自行区分的场景下，你必须使用生命周期注解`。
  当你每次读到一些“必须指定一个生命周期”或者“给这个引用加一个生命周期”这样的句子时，它指的就是生命周期注解。不可能给一个值加一个新的生命周期。

  有两种使用'static 的典型方式。

  - 在一个引用上，作为显式生命周期注解
  - 或者在泛型参数上，作为生命周期的约束(bound)
    &'static 意为该引用在整个程序内是有效的。它指向的这个数据既不会移动也不会改变，它永远都是可用的

- 第 17 天：数组，循环和迭代器

  当能返回迭代器(iterator)自身的时候，那么用 .collect() 返回一个特定的数据类型就是个很差的做法。返回迭代器可以保持灵活性，并且保留了 Rust 开发者期待的懒执行

- 第 18 天：异步

  Rust 的标准库用 `Future` 特质(trait)定义了一个异步(asynchronous)任务的外表模样。但是实现 Future (预期值)，并不足以让任务变成“异步(async)”，你还需要一些别的事物来管理它们。你需要一个装载预期值(futures)的桶，这个桶会检查哪些预期值(futures)已经完成，并通知其等待方。你需要一个执行器(executor)和一个反应器(reactor)，这有点像 node.js 的事件循环

- 第 19 天：开始一个大型项目
- 第 20 天：命令行界面(CLI)参数和日志
- 第 21 天：创建和运行 WebAssembly
- 第 22 天：使用 JSON
- 第 23 天：骗过借用检查器
  使用 Rc 和 Arc 在 Rust 中做引用计数
- 24 天：箱(Crates)和工具
  写 Rust 代码就像跟一个脾气暴躁的老开发人员结对编程，他能发现任何问题。
