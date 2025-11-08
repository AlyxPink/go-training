# Hints for Array Basics

## Level 1: Getting Started

- Arrays in Go have a fixed size that is part of their type
- Use `len(arr)` to get the array length
- Array indices start at 0 and go to `len(arr)-1`
- Accessing out-of-bounds indices causes a runtime panic

## Level 2: Finding Maximum

- Initialize max with the first element
- Iterate through the rest of the array
- Update max when you find a larger value
- Alternative: use the math.MinInt constant for initial value

## Level 3: Array Rotation

- Rotation by k is equivalent to rotation by k % len(arr)
- Right rotation by k = move last k elements to the front
- You can use a temporary array to hold the result
- Think about how indices map: old index i goes to new index (i+k) % len

## Level 4: Finding Duplicates

- Use a map to track which values you've seen
- Use another map or slice to track which are duplicates
- Make sure to return each duplicate only once
- The result can be a slice (variable size)

## Level 5: Edge Cases

Consider:
- What if k is larger than array length? (use modulo)
- What if k is negative? (handle separately or normalize)
- What if the array has no duplicates?
- What about single-element or empty arrays?
