import Foundation

/// Manages offline queue for crashes and events
final class OfflineQueue {
    
    // MARK: - Properties
    
    private let fileManager = FileManager.default
    private let queueDirectory: URL
    private let encoder = JSONEncoder()
    private let decoder = JSONDecoder()
    
    // MARK: - Initialization
    
    init() {
        let documents = fileManager.urls(for: .documentDirectory, in: .userDomainMask).first!
        queueDirectory = documents.appendingPathComponent("BackfeedrQueue", isDirectory: true)
        
        // Create directory if needed
        try? fileManager.createDirectory(at: queueDirectory, withIntermediateDirectories: true)
    }
    
    // MARK: - Queue Operations
    
    /// Add an item to the queue
    func enqueue<T: Encodable>(_ item: T, type: QueueItemType) {
        let filename = "\(type.rawValue)_\(UUID().uuidString).json"
        let url = queueDirectory.appendingPathComponent(filename)
        
        do {
            let data = try encoder.encode(item)
            try data.write(to: url)
        } catch {
            print("[Backfeedr] Failed to enqueue item: \(error)")
        }
    }
    
    /// Get all queued items
    func dequeueAll() -> [QueueItem] {
        guard let files = try? fileManager.contentsOfDirectory(at: queueDirectory, includingPropertiesForKeys: nil) else {
            return []
        }
        
        return files.compactMap { url in
            guard let data = try? Data(contentsOf: url),
                  let item = try? decoder.decode(QueueItem.self, from: data) else {
                return nil
            }
            return item
        }
    }
    
    /// Remove a processed item
    func remove(_ item: QueueItem) {
        let url = queueDirectory.appendingPathComponent(item.filename)
        try? fileManager.removeItem(at: url)
    }
    
    /// Clear all queued items
    func clear() {
        try? fileManager.removeItem(at: queueDirectory)
        try? fileManager.createDirectory(at: queueDirectory, withIntermediateDirectories: true)
    }
}

// MARK: - Types

extension OfflineQueue {
    enum QueueItemType: String {
        case crash = "crash"
        case event = "event"
    }
    
    struct QueueItem: Codable {
        let filename: String
        let type: QueueItemType
        let payload: Data
        let createdAt: Date
    }
}
