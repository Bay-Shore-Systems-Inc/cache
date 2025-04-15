module.exports = {
    header: '# Changelog\n\nAll notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.\n',
    skip: {
        tag: false
    },
    types: [
        { type: 'feat', section: 'Features' },
        { type: 'fix', section: 'Bug Fixes' },
        { type: 'chore', section: 'Chores' },
        { type: 'docs', section: 'Documentation' },
        { type: 'refactor', section: 'Code Refactoring' },
        { type: 'perf', section: 'Performance Improvements' },
        { type: 'test', section: 'Tests' }
    ]
};
