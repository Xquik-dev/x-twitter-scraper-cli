// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

/**
 * @name Code comment
 * @description Code has no comments. A comment needs a legal marker in its file's first comment group, a functional directive or an https link that proves it. A library package may keep 1 doc comment for ST1000.
 * @kind problem
 * @problem.severity error
 * @id go/no-code-comments
 */

import go

predicate proven(Comment comment) {
  exists(CommentGroup group, int index, int later |
    comment = group.getComment(index) and
    later >= index and
    (
      group.getComment(later).getText().regexpMatch("(?s).*https://.*")
      or
      group = comment.getFile().getCommentGroup(0) and
      group.getComment(later).getText().regexpMatch("(?s).*(SPDX-|REUSE-Ignore).*")
    )
  )
}

predicate directive(Comment comment) {
  comment.getText().regexpMatch("(?s)(go:[a-z]+|\\+build|nolint|line )\\b.*|\\s*(Deprecated: .*|(//\\s*)?indirect\\s*)")
}

predicate documentsPackage(GoFile file) {
  exists(file.getDocumentation()) and
  not file.getPackageName() = "main" and
  not file.getBaseName().matches("%\\_test.go")
}

predicate packageDoc(Comment comment) {
  exists(GoFile file |
    documentsPackage(file) and
    comment = file.getDocumentation().getComment(_) and
    not exists(GoFile earlier |
      documentsPackage(earlier) and
      earlier.getParentContainer() = file.getParentContainer() and
      earlier.getPackageName() = file.getPackageName() and
      earlier.getRelativePath() < file.getRelativePath()
    )
  )
}

from Comment comment
where
  exists(comment.getFile().getRelativePath()) and
  not proven(comment) and
  not directive(comment) and
  not packageDoc(comment)
select comment, "Code has no comments. Put the meaning in names, types and tests."
