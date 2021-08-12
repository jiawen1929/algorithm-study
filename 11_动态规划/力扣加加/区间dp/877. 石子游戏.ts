/**
 * @param {number[]} piles
 * @return {boolean}
 * @description
 * 偶数堆石子排成一行，每堆都有正整数颗石子 piles[i] 。
 * 每回合，玩家从行的开始或结束处取走整堆石头。手中石子最多的玩家获胜。
 * 当先手赢得比赛时返回 true ，当后手赢得比赛时返回 false
 * @summary return true
 * 先手取的位置必定限制了后手能取的位置，即： 先手取首位后手只能取偶数位，先手取末位后手只能取奇数位。 因此先手只需要计算好奇偶数位总和，则必赢。
 */
// const stoneGame = function (piles: number[]): boolean {
//   return true
// }
const stoneGame = function (piles: number[]): boolean {
  const len = piles.length
  // dp[i][j]：表示先手玩家（亚历克斯）与后手玩家（李）在区间 [i, j] 之间互相拿，
  // 先手玩家比后手玩家多的最大石子个数。这是个差值，而且是个最大差值。
  const dp = Array.from({ length: len }, () => Array(len).fill(Infinity))
  piles.forEach((pile, index) => (dp[index][index] = pile))
  for (let i = len - 1; i >= 0; i--) {
    for (let j = i + 1; j < len; j++) {
      // 对于先手玩家，有两种拿法：
      dp[i][j] = Math.max(piles[i] - dp[i + 1][j], piles[j] - dp[i][j - 1])
    }
  }
  console.table(dp)
  return dp[0][len - 1] >= 0
}

console.log(stoneGame([5, 3, 4, 5]))

export {}
