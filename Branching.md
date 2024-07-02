## Branching Guidelines 🌿

To keep our workflow organized and efficient, let's follow these branching guidelines:

1. **Main Branch (`main` or `master`)** 🌟

    - This branch should always be in a deployable state. Only thoroughly tested and reviewed code should be merged here.

2. **Development Branch (`dev`)** 🛠️

    - Use this branch for ongoing development. It should contain the latest delivered development changes for the next release. All feature branches should be merged here.

3. **Feature Branches (`feature/author-name/feature-name`)** ✨

    - **Purpose:** Develop new features or user stories.
    - **Process:**
        1. Branch off from `dev`.
        2. Implement the feature, ensuring to follow coding guidelines.
        3. Regularly pull updates from `dev` to keep your branch up to date.
        4. Once complete, create a pull request (PR) to merge back into `dev`.
        5. Address any feedback from the code review.
        6. Merge the PR into `dev`.

4. **Bugfix Branches (`bugfix/author-name/bug-description`)** 🐞

    - **Purpose:** Fix bugs identified in the `dev` or `main` branch.
    - **Process:**
        1. Branch off from `dev` (or `main` if it's a bug in production).
        2. Fix the bug, ensuring to follow coding guidelines.
        3. Test the fix thoroughly.
        4. Create a PR to merge back into `dev` (and `main` if applicable).
        5. Address any feedback from the code review.
        6. Merge the PR into `dev` (and `main` if applicable).

5. **Hotfix Branches (`hotfix/hotfix-description`)** 🔥

    - **Purpose:** Address critical issues that need immediate attention in production.
    - **Process:**
        1. Branch off from `main`.
        2. Implement the hotfix, ensuring to follow coding guidelines.
        3. Test the hotfix thoroughly.
        4. Create a PR to merge back into `main`.
        5. Address any feedback from the code review.
        6. Merge the PR into `main`.
        7. Create another PR to merge the hotfix into `dev` to ensure the fix is included in ongoing development.
        8. Merge the PR into `dev`.

6. **Release Branches (`release/release-version`)** 🚀

    - **Purpose:** Prepare for a new production release.
    - **Process:**
        1. Branch off from `dev` when ready for release.
        2. Perform final testing and make any necessary fixes.
        3. Update version numbers and documentation as needed.
        4. Create a PR to merge back into `main`.
        5. Address any feedback from the code review.
        6. Merge the PR into `main`.
        7. Tag the `main` branch with the new version number.
        8. Create another PR to merge the release branch back into `dev` to keep it up to date.
        9. Merge the PR into `dev`.

7. **Version Tags** 🏷️
    - **Purpose:** Mark specific points in history as being important, typically used for releases.
    - **Tagging convention:** `v1.0.0`

By following these guidelines, we'll ensure our development process remains smooth and organized. Happy branching! 🌳
