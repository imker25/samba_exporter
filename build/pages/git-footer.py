# ######################################################################################
# Copyright 2021 by tobi@backfrak.de. All
# rights reserved. Use of this source code is governed
# by a BSD-style license that can be found in the
# LICENSE file.
# ######################################################################################
# Module called by mkdocs-macros-plugin to create a footer on each page that contains 
# git info of the source '*.md file
# ######################################################################################
import os
import mkdocs_macros

def define_env(env):
    """
    This is the hook for the functions (new form)
    """
    env.variables['cwd'] = os.getcwd()

    
def on_post_page_macros(env):
    """
    Actions to be done after macro interpretation,
    when the markdown is already generated
    """
    gitInfo = mkdocs_macros.context.get_git_info()
    # This will add a (Markdown or HTML) footer
    footer = "\n<br>\n---\n<br><sub>Last change at %s by %s - commit: %s<sub>"  %(gitInfo['date'].strftime("%b %d, %Y %H:%M:%S"), gitInfo['author'], gitInfo['short_commit'])
    env.raw_markdown += footer

