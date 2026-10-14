// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

/**
 * @name Code comment
 * @description Code has no comments. A comment needs a legal marker, a functional directive or an https link that proves it.
 * @kind problem
 * @problem.severity error
 * @id go/no-code-comments
 */

import go

predicate proven(Comment comment) {
  exists(CommentGroup group, int index, int later |
    comment = group.getComment(index) and
    later >= index and
    group.getComment(later).getText().regexpMatch("(?s).*(https://|SPDX-|REUSE-Ignore).*")
  )
}

predicate directive(Comment comment) {
  comment.getText().regexpMatch("(?s)(go:[a-z]+|\\+build|nolint|line )\\b.*|\\s*(Deprecated: .*|(//\\s*)?indirect\\s*)")
}

from Comment comment
where
  exists(comment.getFile().getRelativePath()) and
  not proven(comment) and
  not directive(comment)
select comment, "Code has no comments. Put the meaning in names, types and tests."
