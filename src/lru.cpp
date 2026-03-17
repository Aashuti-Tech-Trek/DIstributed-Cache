/**
 * ============================================================================
 * LRU Cache Implementation - Source
 * ============================================================================
 */

#include "lru.h"
#include <iostream>

/**
 * Constructor: Initialize LRU cache with given capacity
 * 
 * Example:
 *   LRUCache cache(5);  // Can store max 5 items
 */
LRUCache::LRUCache(int cap) {
    capacity = cap;
}

/**
 * GET Operation
 * 
 * Algorithm:
 *   1. Check if key exists in map
 *   2. If YES: Get value, move to front, return value
 *   3. If NO: Return "NULL"
 * 
 * Example:
 *   Initial state: [c, b, a]
 *   cache.get("a")
 *   Result: [a, c, b]  ← 'a' moved to front
 */
string LRUCache::get(const string& key) {
    // Key not found
    if (cacheMap.find(key) == cacheMap.end()) {
        return "NULL";
    }
    
    // Get iterator to element in list
    auto node = cacheMap[key];
    string value = node->second;
    
    // Move to FRONT (most recently used)
    cacheList.erase(node);
    cacheList.push_front({key, value});
    cacheMap[key] = cacheList.begin();
    
    return value;
}

/**
 * PUT Operation (Insert or Update)
 * 
 * Algorithm:
 *   1. If key exists: Remove old entry
 *   2. If cache full: Remove BACK (least recent)
 *   3. Add new entry to FRONT
 *   4. Update map
 * 
 * Example 1 - Update:
 *   cache.put("a", "1")  // 'a' exists, update to "1" and move to front
 * 
 * Example 2 - Eviction:
 *   Capacity = 3, Cache = [c, b, a]
 *   cache.put("d", "4")
 *   Result: [d, c, b]  ← 'a' evicted (least recent)
 */
void LRUCache::put(const string& key, const string& value) {
    
    // Case 1: Key already exists → Remove old entry
    if (cacheMap.find(key) != cacheMap.end()) {
        cacheList.erase(cacheMap[key]);
    }
    
    // Case 2: Cache is full → Evict least recently used (BACK)
    else if ((int)cacheList.size() == capacity) {
        auto lastElement = cacheList.back();
        string lastKey = lastElement.first;
        
        cacheMap.erase(lastKey);      // Remove from map
        cacheList.pop_back();         // Remove from list
    }
    
    // Case 3: Add new entry to FRONT
    cacheList.push_front({key, value});
    cacheMap[key] = cacheList.begin();
}

/**
 * DEBUG: Display cache contents (front to back)
 * Shows current state of cache
 * 
 * Output format:
 *   key1:value1 key2:value2 key3:value3
 * 
 * Example:
 *   [a:1, c:3, b:2]  (front to back)
 */
void LRUCache::display() {
    cout << "Cache contents (front to back): ";
    for (const auto& p : cacheList) {
        cout << p.first << ":" << p.second << " ";
    }
    cout << "\n";
}
