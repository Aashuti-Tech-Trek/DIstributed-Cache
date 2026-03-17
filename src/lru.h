/**
 * ============================================================================
 * LRU Cache Implementation
 * ============================================================================
 * 
 * Purpose:
 * --------
 * Implements a Least Recently Used (LRU) eviction cache using:
 *   - Doubly Linked List: Maintains order of access (most recent at front)
 *   - Hash Map: O(1) lookups
 * 
 * Why LRU?
 * --------
 * - Prevents unlimited memory growth
 * - Evicts least-used items first
 * - Common in real caches (Redis, Memcached, browser caches)
 * 
 * Time Complexity:
 *   - GET: O(1)
 *   - PUT (insert/update): O(1)
 *   - EVICTION: O(1)
 * 
 * Space Complexity: O(capacity)
 * 
 * How It Works:
 * ============
 * 
 * Doubly Linked List:
 *   [Most Recent] ←→ [Recent] ←→ [Old] ←→ [Least Recent]
 *        front()                                back()
 * 
 * When GET: Move accessed item to FRONT
 * When PUT: Add/Update at FRONT
 * When FULL: Remove BACK item
 * 
 * Example:
 * --------
 * LRUCache cache(3);  // capacity = 3
 * 
 * cache.put("a", "1");
 * cache.put("b", "2");
 * cache.put("c", "3");     // [c(3), b(2), a(1)]
 * cache.get("a");          // [a(1), c(3), b(2)]  ← 'a' moves to front
 * cache.put("d", "4");     // [d(4), a(1), c(3)]  ← 'b' evicted (least recent)
 * 
 * ============================================================================
 */

#ifndef LRU_H
#define LRU_H

#include <unordered_map>
#include <list>
#include <string>

using namespace std;

// Represents a Key-Value pair in cache
typedef pair<string, string> CachePair;

// Iterator type for the doubly linked list
typedef list<CachePair>::iterator CacheIterator;

class LRUCache {
private:
    int capacity;
    
    // Doubly linked list: most recently used at FRONT
    // Each node = (key, value)
    list<CachePair> cacheList;
    
    // HashMap: key → position in list
    // Fast lookup of where key is stored
    unordered_map<string, CacheIterator> cacheMap;
    
public:
    /**
     * Constructor
     * @param cap - Maximum number of items cache can hold
     */
    LRUCache(int cap);
    
    /**
     * Retrieve value by key
     * 
     * If key exists:
     *   - Return the value
     *   - Move key to FRONT (mark as recently used)
     * 
     * If key doesn't exist:
     *   - Return "NULL"
     * 
     * @param key - The key to look up
     * @return Value associated with key, or "NULL" if not found
     */
    string get(const string& key);
    
    /**
     * Insert or update key-value pair
     * 
     * Cases:
     *   1. Key exists → Update value, move to FRONT
     *   2. Key new + space available → Add to FRONT
     *   3. Key new + cache full → Evict BACK item, add to FRONT
     * 
     * @param key - The key to insert/update
     * @param value - The value to store
     */
    void put(const string& key, const string& value);
    
    /**
     * Debug: Print cache contents (front to back)
     * Usage: cache.display();
     */
    void display();
};

#endif
