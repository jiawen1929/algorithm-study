/**
 * @param {number[]} nums
 * @return {number}
 * @summary 快指针是读指针， 慢指针是写指针
 */

var removeDuplicates = function (nums) {
  // // 双指针 没见过的就搬过来
  let slowP = 0
  for (let fastP = 0; fastP < nums.length; fastP++) {
    if (nums[fastP] !== nums[slowP]) {
      // 先移后写
      slowP++
      nums[slowP] = nums[fastP]
    }
  }

  // 原地移除 类似于链表的slow.next = null
  nums.length = slowP + 1
  return slowP + 1
}

const a = [0, 0, 1, 1, 1, 2, 2, 3, 3, 4]
console.log(removeDuplicates(a))
console.log(a)
