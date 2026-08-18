# Changelog

## Next release

- Add `SkipRestrictionChecks` to `user.CreateParams`. Creating a user now runs the instance's allowlist, blocklist, blocked disposable email domains and blocked email subaddresses, the same as sign-up does; set this to exempt a single create, for a backend creating a user it already trusts.

- Add support for the OAuth Applications API. Added the oauthapplication package for API operations and a clerk.OAuthApplication type.
- Add support for multiple invitation templates with the `TemplateSlug` field in `invitation.Create`.
- Add support for listing and creating waitlist entries with the `waitlistentry.List` and `waitlistentry.Create` methods.
- Add support for fetching an organization with its members count, via a new `organizations.GetWithParams` method.
- Add support for session reverification with `SessionClaims.NeedsReverification()`, `SessionReverificationPolicy`, predefined policies like `SessionReverificationStrict`, and middleware via `http.NeedsSessionReverification()`.
- Add support for OR-ing audit log list filters via `audit_logs.ListParams.FilterMatch`. Pass `clerk.AuditLogFilterMatchAny` to match events that satisfy any one of `Subject`, `Type`, `Actor`, `TraceID`, `ClientID`, or `ImpersonatorUserID`. Defaults to `clerk.AuditLogFilterMatchAll` (AND) when omitted, preserving existing behavior.
- Add support for marking an organization domain's ownership as verified via the `organizationdomain.VerifyOwnership` method. The `OrganizationDomain` type now exposes `OwnershipVerification` and `AffiliationVerification`; the existing `Verification` field is deprecated in favor of `AffiliationVerification`.
- Add support for sending transactional emails with the new `email` package (`email.Send`). The recipient (`SendParams.To`, a `Recipient`) accepts either an `Address` or a `UserID`; the two are mutually exclusive, and when `UserID` is set Clerk resolves that user's primary email address server-side.
- Add a Directory-named surface for Directory Sync, which also serves non-SCIM providers such as Google Workspace. New `directory` and `directorygrouprolemapping` packages call the `/directories` routes, and new `clerk.Directory`, `clerk.DirectoryList`, `clerk.DirectoryGroup`, and `clerk.DirectoryGroupRoleMapping` types replace their SCIM-named counterparts. `DirectoryGroupRoleMapping.Role` is an `OrganizationRole`, so roles whose permissions are objects deserialize correctly; the deprecated `SCIMGroupRoleMappingRole` cannot decode them. The `scimdirectory` and `scimgrouprolemapping` packages and the `SCIM*` types are deprecated but unchanged, and keep calling `/scim_directories`. Note that the API selects the `object` discriminator by route, so the same Go type reports `scim_directory` from the deprecated package and `directory` from the new one.
- Add `CustomAttribute.DirectoryPath`, the new name for the deprecated `CustomAttribute.SCIMPath`. The API accepts either name on write and returns both on read, and rejects a request that sends both with different values. `CustomAttribute` now marshals exactly one of the two, so read-modify-write keeps working: a decoded path lands in `SCIMPath` with `DirectoryPath` left nil, and on write `DirectoryPath` is sent when set, otherwise `SCIMPath`. Clearing a path still means leaving both nil.

## 2.2.0

- Add support for bulk invitation creation with the `invitation.BulkCreate` method.
- Add `NameQuery` to `user.ListParams`.

## 2.1.1

- Add `EmailAddressQuery`, `PhoneNumberQuery` and `UsernameQuery` to `user.ListParams`.
- Add support for `missing_member_with_elevated_permissions` checks to the `organization.List` method.

## 2.1.0

- Add support for sign in tokens API operations.
- Add `LegalAcceptedAt` to `User` and the ability to `SkipLegalChecks` when creating or updating a `User`.
- Add `EmailAddressQuery`, `PhoneNumberQuery` and `UsernameQuery` to `user.ListParams`.
- Add `RoleName` field to `OrganizationInvitation` and `OrganizationMembership`.
- Add support for deleting a user's external account via the `user.DeleteExternalAccount` method.
- Add support for listing all organization invitations for a user with the `user.ListOrganizationInvitations` method.
- Add support for listing all organization invitations for an instance with the `organizationinvitation.ListFromInstance` method.
- Add `RequestingUserID` parameter to `organizationinvitation.RevokeParams`.
- Update go-jose dependency to v3.0.3.

## 2.0.9

## 2.0.4

- Add `IgnoreDotsForGmailAddresses` field to `InstanceRestrictions` and `instancesettings.UpdateRestrictionsParams` (#293).

## 2.0.0

- Initial version for changelog.
- Complete rewrite, new library API.
