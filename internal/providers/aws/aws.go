package aws

// This file contains type definitions that are used by the AWS provider.
// The implementation of the provider methods has been moved to use the adapter pattern,
// delegating to the cloud layer through UI operation interfaces.

// Note: The redundant type definitions have been removed in favor of using
// the cloud types directly. All code should now use the types defined in
// internal/cloud/types.go instead.
