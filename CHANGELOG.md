# Changelog

## Next release

## 2.3.1

- Fix parsing the org claims in the new JWT session token version.

## 2.3.0

- Add support for the OAuth Applications API. Added the oauthapplication package for API operations and a clerk.OAuthApplication type.
- Add support for multiple invitation templates with the `TemplateSlug` field in `invitation.Create`.
- Add support for listing and creating waitlist entries with the `waitlistentry.List` and `waitlistentry.Create` methods.
- Add support for fetching an organization with its members count, via a new `organizations.GetWithParams` method.
- Add support for session reverification with `SessionClaims.NeedsReverification()`, `SessionReverificationPolicy`, predefined policies like `SessionReverificationStrict`, and middleware via `http.NeedsSessionReverification()`.
- Handles new JWT session token version

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
