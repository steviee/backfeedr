import Foundation

/// Main entry point for the Backfeedr SDK
public final class Backfeedr {
    
    // MARK: - Properties
    
    /// Shared instance (singleton)
    public static let shared = Backfeedr()
    
    /// Configuration
    private var configuration: Configuration?
    
    /// Crash reporter
    private(set) lazy var crashReporter = CrashReporter(configuration: configuration)
    
    /// Event tracker
    private(set) lazy var eventTracker = EventTracker(configuration: configuration)
    
    /// Private queue for background operations
    private let queue = DispatchQueue(label: "dev.backfeedr.sdk", qos: .utility)
    
    // MARK: - Initialization
    
    private init() {}
    
    // MARK: - Configuration
    
    /// Configure the SDK with endpoint and API key
    /// - Parameters:
    ///   - endpoint: The backfeedr server URL (e.g., "https://crashes.example.com")
    ///   - apiKey: The API key (e.g., "bf_live_...")
    ///   - config: Optional additional configuration
    public func configure(
        endpoint: String,
        apiKey: String,
        config: Configuration = .init()
    ) {
        self.configuration = Configuration(
            endpoint: endpoint,
            apiKey: apiKey,
            scrubPII: config.scrubPII,
            enableHMAC: config.enableHMAC,
            sessionTimeout: config.sessionTimeout
        )
        
        // Initialize components
        crashReporter = CrashReporter(configuration: self.configuration)
        eventTracker = EventTracker(configuration: self.configuration)
        
        // Start session tracking
        eventTracker.startSession()
        
        // Upload any queued events
        queue.async { [weak self] in
            self?.flushQueue()
        }
    }
    
    // MARK: - Crash Reporting
    
    /// Capture a non-fatal error
    /// - Parameters:
    ///   - error: The error to report
    ///   - context: Additional context (will be scrubbed for PII)
    public func capture(_ error: Error, context: [String: Any]? = nil) {
        guard configuration != nil else {
            print("[Backfeedr] Warning: SDK not configured")
            return
        }
        crashReporter.report(error: error, context: context)
    }
    
    // MARK: - Event Tracking
    
    /// Track a custom event
    /// - Parameters:
    ///   - name: Event name
    ///   - properties: Event properties (will be scrubbed for PII)
    public func track(_ name: String, properties: [String: Any]? = nil) {
        guard configuration != nil else {
            print("[Backfeedr] Warning: SDK not configured")
            return
        }
        eventTracker.track(name: name, properties: properties)
    }
    
    /// Leave a breadcrumb for crash context
    /// - Parameter message: The breadcrumb message
    public func breadcrumb(_ message: String) {
        crashReporter.addBreadcrumb(message)
    }
    
    // MARK: - Queue Management
    
    /// Flush the offline queue immediately
    public func flushQueue() {
        guard configuration != nil else { return }
        // TODO: Implement queue flush
    }
}

// MARK: - Configuration

extension Backfeedr {
    public struct Configuration {
        let endpoint: String
        let apiKey: String
        let scrubPII: Bool
        let enableHMAC: Bool
        let sessionTimeout: TimeInterval
        
        public init(
            endpoint: String = "",
            apiKey: String = "",
            scrubPII: Bool = true,
            enableHMAC: Bool = true,
            sessionTimeout: TimeInterval = 300
        ) {
            self.endpoint = endpoint
            self.apiKey = apiKey
            self.scrubPII = scrubPII
            self.enableHMAC = enableHMAC
            self.sessionTimeout = sessionTimeout
        }
    }
}
