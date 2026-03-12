// swift-tools-version:5.9
import PackageDescription

let package = Package(
    name: "BackfeedrKit",
    platforms: [
        .iOS(.v15),
        .macOS(.v12),
        .watchOS(.v8),
        .tvOS(.v15),
        .visionOS(.v1)
    ],
    products: [
        .library(
            name: "BackfeedrKit",
            targets: ["BackfeedrKit"]
        ),
    ],
    dependencies: [
        // No external dependencies for MVP
    ],
    targets: [
        .target(
            name: "BackfeedrKit",
            path: "Sources/BackfeedrKit"
        ),
        .testTarget(
            name: "BackfeedrKitTests",
            dependencies: ["BackfeedrKit"],
            path: "Tests/BackfeedrKitTests"
        ),
    ]
)
