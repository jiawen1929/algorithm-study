# TypeScript 实现 StampedLock

StampedLock 是 Java 8 中引入的一种高性能锁，它支持三种模式：写锁、读锁和乐观读。与 ReentrantReadWriteLock 不同，StampedLock 不是可重入的，但它提供了**乐观读取**功能，在读多写少的场景下有更好的性能。

## StampedLock 实现详解

StampedLock 是一种高级锁实现，提供了三种访问模式：写锁、读锁和乐观读。它不是可重入的，但在读多写少的场景下性能优于 ReentrantReadWriteLock。

### 核心概念

1. **戳记 (Stamp)**:

   - 每个锁操作都会返回一个非零的戳记值
   - 戳记用于验证锁状态和进行解锁操作

2. **锁模式**:
   - **写锁**: 排他模式，获取锁时阻止所有其他线程获取任何形式的锁
   - **读锁**: 共享模式，允许多个线程同时获取读锁
   - **乐观读**: 无锁模式，获取当前戳记但不实际锁定资源，后续可验证数据是否变化

### 关键特性

1. **戳记验证**:

   - 使用戳记而不是布尔值来表示锁状态
   - 戳记允许检测数据自上次访问后是否已更改

2. **乐观读**:

   - 允许在不获取锁的情况下读取数据
   - 通过 `validate()` 方法检查读取是否有效
   - 适用于读操作远多于写操作的场景

3. **锁转换**:

   - 支持锁升级 (读锁→写锁)
   - 支持锁降级 (写锁→读锁)
   - 支持乐观读转换为读/写锁

4. **非阻塞尝试**:
   - 所有锁操作都有非阻塞的尝试版本
   - 允许在无法立即获取锁时执行替代逻辑

### 实现细节

1. **内部状态**:

   - `_state`: 锁的当前状态 (未锁定/写锁定/读锁定)
   - `_writeStamp`: 当前写锁的戳记
   - `_readCount`: 当前读锁的数量
   - `_stamp`: 全局戳记计数器

2. **等待队列**:

   - `_writeWaiters`: 等待获取写锁的线程队列
   - `_readWaiters`: 等待获取读锁的线程队列

3. **优先策略**:

   - 写锁优先，防止写饥饿
   - 当资源可用时，先唤醒写线程，没有写线程时才唤醒读线程

4. **便利方法**:
   - `withReadLock()`: 获取读锁执行操作，自动释放
   - `withWriteLock()`: 获取写锁执行操作，自动释放
   - `withOptimisticRead()`: 使用乐观读执行操作，失败时回退到读锁

### 使用场景

1. **读多写少的场景**:

   - 数据缓存
   - 配置信息访问

2. **需要低延迟读取的场景**:

   - 乐观读提供无阻塞的快速读取
   - 适合大多数读取不会与写入冲突的情况

3. **读写锁升级/降级需求**:

   - 开始时只需读取，后续可能需要修改
   - 先需要独占修改，后续只需共享读取

4. **高并发系统**:
   - 减少锁竞争
   - 提高并发读取的吞吐量

这个实现提供了Java StampedLock的所有核心功能，并针对TypeScript的异步特性进行了适配。通过使用戳记和三种锁模式，它为各种并发访问场景提供了灵活和高效的解决方案。
