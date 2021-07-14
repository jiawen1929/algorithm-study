"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
class Node {
    constructor(value, next) {
        this.value = value;
        this.next = next;
    }
}
const a = new Node(1);
const b = new Node(2);
const c = new Node(3);
a.next = b;
b.next = c;
const addTwo = (l1, l2) => {
    // l3是需要返回的链表
    const l3 = new Node(0);
    let n1 = l1;
    let n2 = l2;
    let n3 = l3;
    // 进位
    let overflow = 0;
    // 遍历链表
    while (n1 || n2) {
        const n1Value = n1?.value || 0;
        const n2Value = n2?.value || 0;
        const sum = n1Value + n2Value + overflow;
        overflow = Math.floor(sum / 10);
        n3.next = new Node(sum % 10);
        // 移动所有节点
        n1 = n1?.next;
        n2 = n2?.next;
        n3 = n3.next;
    }
    // 最后一个进位
    if (overflow) {
        n3.next = new Node(overflow);
    }
    return l3;
};
console.log(addTwo(a, a));
// O(n)
