const { execSync } = require("child_process");

const args = process.argv.slice(2);
const forceRelease = args.includes("--ignore-empty"); // This will let you create a new release even if the changelog would be empty

try {
    console.log("Checking for release-worthy commits...");
    const output = execSync("npx standard-version --dry-run", {
        encoding: "utf-8",
        stdio: ["pipe", "pipe", "pipe"],
    });

    const hasReleaseContent = /### (Features|Bug Fixes|Performance Improvements|Breaking Changes)/i.test(output);

    if (!hasReleaseContent && !forceRelease) {
        console.log("No release-worthy commits found. Skipping version bump.")
        process.exit(0);
    }

    if (forceRelease) {
        console.log("Forcing release despite no changelog content (--ignore-empty)...")
    } else {
        console.log("Commits detected proceeding with creation of new release...");
    }

    execSync("npx standard-version", { stdio: "inherit" });
} catch (err) {
    console.error("Error creating release: ", err.message)
    process.exit(1);
}
