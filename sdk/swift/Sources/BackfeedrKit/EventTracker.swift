import Foundation

/// Tracks events and sessions
final class EventTracker {
    
    // MARK: - Properties
    
    private let configuration: Backfeedr.Configuration?
    private let queue = DispatchQueue(label: "dev.backfeedr.events", qos: .utility)
    private var currentSessionID: String?
    private var sessionStartTime: Date?
    
    // MARK: - Initialization
    
    init(configuration: Backfeedr.Configuration?) {
        self.configuration = configuration
    }
    
    // MARK: - Session Management
    
    func startSession() {
        currentSessionID = generateSessionID()
        sessionStartTime = Date()
        
        // Send session_start event
        sendEvent(
            type: .sessionStart,
            name: nil,
            properties: nil
        )
    }
    
    func endSession() {
        guard let startTime = sessionStartTime else { return }
        
        let duration = Date().timeIntervalSince(startTime)
        
        sendEvent(
            type: .sessionEnd,
            name: nil,
            properties: ["duration": duration]
        )
        
        currentSessionID = nil
        sessionStartTime = nil
    }
    
    // MARK: - Event Tracking
    
    func track(name: String, properties: [String: Any]?) {
        sendEvent(
            type: .custom,
            name: name,
            properties: properties
        )
    }
    
    func trackError(_ error: Error, properties: [String: Any]?) {
        var props = properties ?? [:]
        props["error_type"] = String(describing: type(of: error))
        props["error_description"] = error.localizedDescription
        
        sendEvent(
            type: .error,
            name: "error",
            properties: props
        )
    }
    
    // MARK: - Private Methods
    
    private func sendEvent(type: EventType, name: String?, properties: [String: Any]?) {
        guard configuration != nil else { return }
        
        queue.async { [weak self] in
            // TODO: Implement event queuing and sending
            print("[Backfeedr] Event: \(type.rawValue)\(name != nil ? " - \\(name!)" : "")")
        }
    }
    
    private func generateSessionID() -> String {
        UUID().uuidString.replacingOccurrences(of: "-", with: "")
    }
}

// MARK: - Event Types

extension EventTracker {
    enum EventType: String {
        case sessionStart = "session_start"
        case sessionEnd = "session_end"
        case error = "error"
        case custom = "custom"
    }
}
