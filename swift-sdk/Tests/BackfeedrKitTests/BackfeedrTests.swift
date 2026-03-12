import XCTest
@testable import BackfeedrKit

final class BackfeedrTests: XCTestCase {
    
    override func setUp() {
        super.setUp()
    }
    
    override func tearDown() {
        super.tearDown()
    }
    
    func testConfiguration() {
        // Given
        let config = Backfeedr.Configuration(
            endpoint: "https://test.example.com",
            apiKey: "bf_live_test123",
            scrubPII: true,
            enableHMAC: true
        )
        
        // Then
        XCTAssertEqual(config.endpoint, "https://test.example.com")
        XCTAssertEqual(config.apiKey, "bf_live_test123")
        XCTAssertTrue(config.scrubPII)
        XCTAssertTrue(config.enableHMAC)
    }
    
    func testHMACSigning() {
        // Given
        let signer = HMACSigner(apiKey: "bf_live_testkey")
        let timestamp = "2026-03-12T14:22:31.123Z"
        let body = "{\"test\":\"data\"}".data(using: .utf8)!
        
        // When
        let signature = signer.sign(timestamp: timestamp, body: body)
        
        // Then
        XCTAssertEqual(signature.count, 64) // SHA256 hex = 64 chars
        XCTAssertTrue(signature.range(of: "[^a-f0-9]", options: .regularExpression) == nil)
    }
    
    func testPIIScrubbing_email() {
        // Given
        let scrubber = PIIScrubber()
        let input = "Contact us at support@example.com for help"
        
        // When
        let result = scrubber.scrub(input)
        
        // Then
        XCTAssertEqual(result, "Contact us at [REDACTED_EMAIL] for help")
    }
    
    func testPIIScrubbing_phone() {
        // Given
        let scrubber = PIIScrubber()
        let input = "Call +1-555-123-4567"
        
        // When
        let result = scrubber.scrub(input)
        
        // Then
        XCTAssertEqual(result, "Call [REDACTED_PHONE]")
    }
    
    func testPIIScrubbing_dictionary() {
        // Given
        let scrubber = PIIScrubber()
        let input: [String: Any] = [
            "email": "user@example.com",
            "screen": "settings"
        ]
        
        // When
        let result = scrubber.scrub(input)
        
        // Then
        XCTAssertEqual(result["email"] as? String, "[REDACTED_EMAIL]")
        XCTAssertEqual(result["screen"] as? String, "settings")
    }
}
