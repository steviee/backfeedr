import Foundation

/// Scrubs PII from payloads before sending
final class PIIScrubber {
    
    // MARK: - Patterns
    
    /// Regex patterns for detecting PII
    private let patterns: [(pattern: String, replacement: String)] = [
        // Email addresses
        ("[A-Z0-9._%+-]+@[A-Z0-9.-]+\\.[A-Z]{2,}", "[REDACTED_EMAIL]"),
        // Phone numbers (basic)
        ("\\+?[0-9]{1,3}[-. ]?\\(?[0-9]{3}\\)?[-. ]?[0-9]{3}[-. ]?[0-9]{4}", "[REDACTED_PHONE]"),
        // IPv4 addresses
        ("\\b(?:[0-9]{1,3}\\.){3}[0-9]{1,3}\\b", "[REDACTED_IP]"),
        // URLs with query params
        ("https?://[^\\s]+\\?[^\\s]+", "[REDACTED_URL]"),
    ]
    
    // MARK: - Scrubbing
    
    /// Scrub a string value
    func scrub(_ value: String) -> String {
        var result = value
        
        for (pattern, replacement) in patterns {
            if let regex = try? NSRegularExpression(pattern: pattern, options: .caseInsensitive) {
                let range = NSRange(location: 0, length: result.utf16.count)
                result = regex.stringByReplacingMatches(in: result, options: [], range: range, withTemplate: replacement)
            }
        }
        
        return result
    }
    
    /// Scrub a dictionary recursively
    func scrub(_ dictionary: [String: Any]) -> [String: Any] {
        var result: [String: Any] = [:]
        
        for (key, value) in dictionary {
            // Skip known safe keys
            if isSafeKey(key) {
                result[key] = value
                continue
            }
            
            switch value {
            case let str as String:
                result[key] = scrub(str)
            case let dict as [String: Any]:
                result[key] = scrub(dict)
            case let array as [Any]:
                result[key] = scrub(array)
            default:
                result[key] = value
            }
        }
        
        return result
    }
    
    /// Scrub an array recursively
    func scrub(_ array: [Any]) -> [Any] {
        return array.map { value in
            switch value {
            case let str as String:
                return scrub(str)
            case let dict as [String: Any]:
                return scrub(dict)
            case let arr as [Any]:
                return scrub(arr)
            default:
                return value
            }
        }
    }
    
    // MARK: - Helpers
    
    private func isSafeKey(_ key: String) -> Bool {
        let safeKeys = [
            "screen_name",
            "view_controller",
            "button_name",
            "feature_flag",
            "experiment_id",
        ]
        return safeKeys.contains(key.lowercased())
    }
}
