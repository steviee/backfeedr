import Foundation
import CryptoKit

/// Handles HMAC request signing
final class HMACSigner {
    
    // MARK: - Properties
    
    private let apiKey: String
    
    // MARK: - Initialization
    
    init(apiKey: String) {
        self.apiKey = apiKey
    }
    
    // MARK: - Signing
    
    /// Sign a request payload
    /// - Parameters:
    ///   - timestamp: ISO8601 timestamp
    ///   - body: Request body data
    /// - Returns: HMAC signature string
    func sign(timestamp: String, body: Data) -> String {
        // Calculate body hash
        let bodyHash = SHA256.hash(data: body)
        let bodyHashHex = bodyHash.compactMap { String(format: "%02x", $0) }.joined()
        
        // Build payload: timestamp.body_hash
        let payload = "\(timestamp).\(bodyHashHex)"
        
        // Calculate HMAC
        let key = SymmetricKey(data: Data(apiKey.utf8))
        let signature = HMAC<SHA256>.authenticationCode(for: Data(payload.utf8), using: key)
        
        return signature.compactMap { String(format: "%02x", $0) }.joined()
    }
    
    /// Generate current timestamp in ISO8601 format
    func timestamp() -> String {
        let formatter = ISO8601DateFormatter()
        formatter.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
        return formatter.string(from: Date())
    }
}

// MARK: - Request Headers

extension HMACSigner {
    struct SignedHeaders {
        let timestamp: String
        let signature: String
        let apiKey: String
        
        func asDictionary() -> [String: String] {
            [
                "X-Backfeedr-Key": apiKey,
                "X-Backfeedr-Timestamp": timestamp,
                "X-Backfeedr-Signature": "sha256=\(signature)"
            ]
        }
    }
    
    /// Sign a request and return headers
    func signRequest(body: Data) -> SignedHeaders {
        let ts = timestamp()
        let sig = sign(timestamp: ts, body: body)
        
        return SignedHeaders(
            timestamp: ts,
            signature: sig,
            apiKey: apiKey
        )
    }
}
