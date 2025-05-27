# Algorithms in Python (modified from the excellent: Grokking Algorithms) - see also http://www.integralist.co.uk/posts/bigo.html for details on understanding Big O notation

[View original Gist on GitHub](https://gist.github.com/Integralist/9763bded76e7d826535a3caeafc3bdff)

## 0. Algorithms in Python.md

- [Binary Search](#file-1-binary-search-py)
- [Selection Sort](#file-2-selection-sort-py)
- [Quick Sort (first index)](#file-3-quick-sort-first-index-py)
- [Quick Sort (middle index)](#file-4-quick-sort-middle-index-py)
- [Quick Sort (random index)](#file-5-quick-sort-random-index-py)
- [Breadth First Search](#file-6-breadth-first-search-py)
- [Dijkstra's Algorithm](#file-7-dijkstra-s-algorithm-py)

## 1. binary search.py

```python
'''
Binary Search algorithm

Description:
    Locate an item within a collection by using a 
    divide and conquer approach.

    In essence, pick the middle of the collection
    then verify if the value is too high or low
    and then reduce the sliding window, now pick 
    the middle again - rinse/repeat

Performance:
    O(Log₂ n)
    Logarithmic time

Meaning:
    12 item collection will take (worst case) ~4 steps
    to locate the specified index using binary search

Note: collection must be sorted
'''

collection = [1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 20, 21]  # 12 items


def binary_search(collection, item):
    print("we're looking for:", item, "\n")
    print("collection:", collection, "\n")

    start = 0
    stop = len(collection) - 1

    while start <= stop:
        middle = round((start + stop) / 2)
        guess = collection[middle]

        print("start: {}\nstop: {}\nmiddle: {}\nguess: {}\n".format(start,
                                                                    stop,
                                                                    middle,
                                                                    guess))

        if guess == item:
            return middle
        if guess > item:
            stop = middle - 1
        else:
            start = middle + 1

    return None


print(binary_search(collection, 9))  # found at index: 4
```

## 2. selection sort.py

```python
'''
Selection Sort algorithm

Description:
    Sorts an unordered list by looping over the list
    n number of times and for each loop identifying
    either the smallest or largest element (which one 
    depends on how you're hoping to sort your list:
    do you want ascending or descending order)

    You'll end up constructing a new ordered list

Performance:
    O(n x n)
    O(n₂)

    Both forms are equivalent.
    You effectively loop a number of times to match collection length
    Then you loop the collection looking for smallest/largest
    Then you mutate the collection so it's smaller by one
    
    So although you're looping over the collection multiple times, 
    you are in fact looping over a slightly smaller collection each time

    Although in reality the notation should be:
        O(n x n - 1)

    But that is invalid Big O notation, so we use
    the O(n₂) or O(n x n) instead
    
    If your collection is 10 items long then this calculates as:
        10 * 10 (i.e. 10 to the power of 2) = 100 operations
        0.1 * 100 = 10 seconds
    
Note: to test algorithm in a different order, the
      quickest change is to turn < into > within 
      the find_smallest function. 
      This would cause program to return [9, 5, 3, 1]
'''

def find_smallest(arr):
    smallest = arr[0]
    smallest_index = 0

    for i, _ in enumerate(arr):
        if arr[i] < smallest:
            smallest = arr[i]
            smallest_index = i

    return smallest_index

def selection_sort(arr):
    temp = []

    for i in range(len(arr)):
        smallest = find_smallest(arr)
        temp.append(arr.pop(smallest))

    return temp

print(selection_sort([5, 9, 3, 1]))  # [1, 3, 5, 9]
```

## 3. quick sort (first index).py

```python
'''
Quick Sort algorithm

Description:
    Sorts an unordered list using recusion and specifically
    using the D&C (Divide and Conquer) approach to problem solving

Explanation:
    You pick a 'pivot' (a random array index)
    You loop the array storing items less than the pivot
    You loop the array storing items greater than the pivot
    Assume less and greater are already sorted
    You can now return 'less' + pivot + 'greater'

    In reality you'll use recursion to then sort both the
    less and greater arrays using the same algorithm

Performance:
    Worst case:
        O(n x n)
        O(n₂)
        
        If collection is 10 items in length (10 * 10 = 100 operations)
        
    Best case:
        O(n Log₂ n)
        
        If collection is 10 items in length (10 * 4 (Log₂10 == 2*2*2*2) = 40 operations)

    Quick Sort's performance depends on the pivot you choose

    In our example code we always pick the zero index as it
    makes the code simpler, but not necessarily the most performant
    
Notes:
    The following implementation doesn't pick a pivot at random
    It instead grabs the first index
'''

def quicksort(arr):
    if len(arr) < 2:
        return arr
    else:
        pivot = arr[0]
        less = [i for i in arr[1:] if i <= pivot]
        greater = [i for i in arr[1:] if i > pivot]
        return quicksort(less) + [pivot] + quicksort(greater)

print(quicksort([10, 5, 2, 3, 7, 0, 9, 12]))  # [0, 2, 3, 5, 7, 9, 10, 12]
```

## 4. quick sort (middle index).py

```python
'''    
Notes:
    The following implementation doesn't pick a pivot at random
    It instead grabs the middle index
'''

items = [10, 5, 2, 3, 7, 0, 9, 12]
print("Original:", items)

def quicksort(arr):
    if len(arr) < 2:
        return arr
    else:
        middle = round(len(arr) / 2)  # grab the middle index
        pivot = arr.pop(middle)
        less = [i for i in arr if i <= pivot]
        greater = [i for i in arr if i > pivot]
        return quicksort(less) + [pivot] + quicksort(greater)

print("Sorted:  ", quicksort(items))

"""
Original: [10, 5, 2, 3, 7, 0, 9, 12]
Sorted:   [0, 2, 3, 5, 7, 9, 10, 12]
"""
```

## 5. quick sort (random index).py

```python
'''    
Notes:
    The following implementation picks a pivot at random
    As it should be, per the original definition of the algorithm's design
'''

from random import randrange

items = [10, 5, 2, 3, 7, 0, 9, 12]
print("Original:", items)

def quicksort(arr):
    if len(arr) < 2:
        return arr
    else:
        middle = randrange(0, len(arr))  # grab a random index
        pivot = arr.pop(middle)
        less = [i for i in arr if i <= pivot]
        greater = [i for i in arr if i > pivot]
        return quicksort(less) + [pivot] + quicksort(greater)

print("Sorted:  ", quicksort(items))

"""
Original: [10, 5, 2, 3, 7, 0, 9, 12]
Sorted:   [0, 2, 3, 5, 7, 9, 10, 12]
"""
```

## 6. Breadth First Search.py

```python
'''
Breadth First Search

Description:
    Breadth-first search tells you if there’s a path from A to B.
    If there’s a path, breadth-first search will find the shortest path.
    A directed graph has arrows, and the relationship follows the direction of the arrow.
    Undirected graphs don’t have arrows, and the relationship goes both ways.

Explanation:
    A graph is made up of 'nodes' and 'edges'.
    In this algorithm there are no 'weights' to the edges (see: Dijkstra’s algorithm)
    An directed graph can also be a 'Tree', but an undirected graph cannot.
    The closest nodes are our 1st 'degree', outer nodes are 2nd, 3rd degrees etc.

    We keep a queue of the 1st degree nodes.
    Pop the first node off the queue and check if this gives you your 'match'

        'match' being the logic that determines if you can stop searching

    If it doesn't match, then add all the 'neighbours' for that node to the queue.
    Rinse and Repeat (e.g. pop next 1st degree node, check it, add neighbours...)
    If the queue is empty, then there were no matches within your network graph.

Performance:
    O(V+E)

    It means V for vertices and E for edges.
    That is, worst case, you have to search the entire graph & follow every edge.
    It also includes O(1) for adding new degrees to the queue

Notes:
    Only provides you with the 'shortest' route, not the 'fastest'
    For 'fastest' route (i.e. weighted graphs) see Dijkstra’s algorithm

    You need to check people in the order they were added to the search list,
    so the search list needs to be a queue.
    Otherwise, you won’t get the shortest path.

    You also need to make sure not to check the same node twice.
    Some nodes from a degree could share a node in the next outer degree.
    It doesn't break the algorithm but is a wasteful operation performance-wise.
'''

graph = {}
graph['you'] = ['alice', 'bob', 'claire']
graph['bob'] = ['anuj', 'peggy']
graph['alice'] = ['peggy']
graph['claire'] = ['thom', 'johnny']
graph['anuj'] = []
graph['peggy'] = []
graph['thom'] = []
graph['jonny'] = []

'''
^^ is your network graph in basic code.
You are in the center and you have 3 neighbours.
We declare the neighbours for everyone in our data model (dict/hash table).
We're looking for someone whose name has 'm' as the last letter.
'''

from collections import deque


def check(key):
    return key[-1] == 'm'


def search(starting_point):
    queue = deque()
    queue += graph[starting_point]  # add starting_point's neighbours initially
    searched = []  # track items we've already checked

    while queue:
        item = queue.popleft()
        if item not in searched:
            if check(item):
                print('found a match: {}'.format(item))
                return True  # we're finished!
            else:
                queue += graph[item]  # add this item's neighbours
                searched.append(item)  # tracked as being checked already

    return False

search('you')  # found a match: thom
```

## 7. Dijkstra's Algorithm.py

```python
'''
Dijkstra's Algorithm

Description:
    Tells you the quickest path from A to B within a weighted graph.

Explanation:
    Create a model of your graph which indicates the following:

        Parent
        Node
        Cost

    Get the node closest to the 'start'.
    Check its neighbours and update the 'Cost' column.
    If a neighbour's cost is updated, then update the parent's too.
    Mark this node as 'processed' so you don't end up processing it again.

    For more details on graph terminology see 'Breadth First Search'.

Performance:
    O(V+E)

    It means V for vertices and E for edges.
    That is, worst case, you have to search the entire graph & follow every edge.

Notes:
    It doesn't work with 'negative' weighted graphs.

Weighted/Directed Graph Example:

         ---> A --->
        |     ^     |
        6     |     1
        |     |     |
    start     3     |-->end
        |     |     |
        2     |     5
        |     |     |
         ---> B --->

The various routes and their costs are:

    * start -> A -> end: 7
    * start -> B -> end: 7
    * start -> B -> A -> end: 6

We can see the last route has more steps but is faster in terms of 'weight'
'''

#################################

'''
Let's start to model our graph...
'''

graph = {}

# Start has two neighbours 'a' and 'b'
# So we model that as follows, with their associated edge weights
graph['start'] = {}
graph['start']['a'] = 6
graph['start']['b'] = 2

# A has one neighbour (end)
graph['a'] = {}
graph['a']['end'] = 1

# B has two neighbours (A and end)
graph['b'] = {}
graph['b']['a'] = 3
graph['b']['end'] = 5

# End has no neighbours
graph['end'] = {}

'''
Now we need a 'costs' hash table...
'''

costs = {}
costs['a'] = 6
costs['b'] = 2
costs['end'] = float('inf')  # set to infinity until we know the cost to reach

'''
Next we need a 'parents' hash table...
'''

parents = {}
parents['a'] = 'start'
parents['b'] = 'start'
parents['end'] = None  # doesn't have one yet until we choose either 'a' or 'b'

'''
Finally we need to track which nodes have already been processed...
'''

processed = []

'''
Oh, and one more array for the purpose of displaying the final route.
It's not used as part of the core algorithm, but utilised by the 'display_route' function
'''

route = []

'''
Here is the implementation...
'''


def find_fastest_path():
    node = find_lowest_cost_node(costs)

    while node is not None:
        cost = costs[node]
        neighbours = graph[node]

        for n in neighbours.keys():
            new_cost = cost + neighbours[n]

            if costs[n] > new_cost:
                costs[n] = new_cost
                parents[n] = node

        processed.append(node)
        node = find_lowest_cost_node(costs)


def find_lowest_cost_node(costs):
    lowest_cost = float('inf')
    lowest_cost_node = None

    for node in costs:
        cost = costs[node]
        if cost < lowest_cost and node not in processed:
            lowest_cost = cost
            lowest_cost_node = node

    return lowest_cost_node


def display_route(node=None):
    # These vars could be passed into the function
    # if you wanted dynamic processing...

    start_node = 'start'  # we look for this 'value'
    end_node = 'end'  # we look for this 'key'

    if node == start_node:  # we're done
        route.append(node)
        reverse_route = list(reversed(route))
        print('Fastest Route: ' + ' -> '.join(reverse_route))
    elif node is None:  # begin processing (this condition only passes once)
        end_parent = parents[end_node]  # get the parent node for the end node
        route.append(end_node)
        display_route(end_parent)
    else:
        parent = parents[node]
        route.append(node)
        display_route(parent)

find_fastest_path()  # mutates global 'costs' & 'parents' arrays
display_route()  # Fastest Route: start -> b -> a -> end
```

