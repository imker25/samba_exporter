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
    footer = "\n---\n<br><sub>Last change at %s by %s - commit: %s<sub>"  %(gitInfo['date'].strftime("%b %d, %Y %H:%M:%S"), gitInfo['author'], gitInfo['short_commit'])
    env.raw_markdown += footer

