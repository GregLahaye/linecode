class Solution:
    def twoSum(self, nums: List[int], target: int) -> List[int]:
        d = {}
        for index, value in enumerate(nums):
            if value not in d:
                d[value] = []

            d[value].append(index)

        for i in range(len(nums)):
            needed = target - nums[i]
            if needed in d:
                a = [x for x in d[needed] if x != i]

                if len(a):
                    return [i, a[0]]