"use strict";
// 思路：
// 找出所有不包含重复字符的字串，返回长度最大的字串
// 滑动窗口
Object.defineProperty(exports, "__esModule", { value: true });
const getLongestSubstring = (str) => {
    let leftPoint = 0;
    let maxLength = 0;
    let tmpMap = new Map();
    for (let rightPoint = 0; rightPoint < str.length; rightPoint++) {
        const element = str[rightPoint];
        // 遇到重复字符，左指针移动到重复字符的下一位
        if (tmpMap.has(element)) {
            leftPoint = tmpMap.get(element) + 1;
        }
        tmpMap.set(element, rightPoint);
        maxLength = Math.max(maxLength, rightPoint - leftPoint + 1);
    }
    return maxLength;
};
console.log(getLongestSubstring('abbcdea'));
