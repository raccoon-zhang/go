import re
import subprocess

podNameFmt = r'pod\s+["\']([^"\']+?)["\'](?:,\s*["\']([^"\']+?)["\'])?'
podGitFmt = r'pod\s+["\']([^"\']+)[\'"]\s*,\s*:git\s*=>\s*["\']([^"\']+)[\'"]'


def findVersion(name):
    try:
        result = subprocess.check_output("pod search %s" % name, shell=True).decode("utf-8")
    except subprocess.CalledProcessError as e:
        return None

    if result == "":
        return None
    
    pattern = r"pod '%s', '~> ([\w.-]+)'" % name
    match = re.search(pattern, result)

    if match == None:
        return None
    else:
        return match.group(1)
    
def findTag(git):
    try:
        result = subprocess.check_output("git ls-remote --tags %s" % git, shell=True).decode("utf-8").split("\n")
    except subprocess.CalledProcessError as e:
        return None

    if result == "":
        return None
    
    pattern = r'refs/tags/([\w.-]+)$'
    tags = []
    for item in result:
        match = re.findall(pattern, item)
        if len(match) != 0:
            tags.append(match[0])

    if len(tags) == 0:
        return None
    else:
        return tags[-1]


def handlLine(line):
    lineContent = line.strip()
    splitContentByComment = line.split('#')
    cleanContent = splitContentByComment[0].rstrip()
    if len(splitContentByComment) > 1:
        commentContent = " #" + splitContentByComment[1]
    else:
        commentContent = ""
    
    # 跳过注释行
    if re.match(r'^\s*#', lineContent) != None:
        return line

    # 跳过非pod行
    if re.match(r'^pod\s+', lineContent) == None:
        return line

    #指定分支和提交的行跳过
    if re.search(r':commit|:branch', cleanContent): # 使用正则表达式搜索 '#' 之前的部分是否包含 ':commit' 或 ':branch'
        return line  

    if re.search(podGitFmt,cleanContent) != None:
        return handlePodGit(cleanContent) + commentContent
    #版本号管理分支
    if re.search(podNameFmt,cleanContent) != None:
        return handlePodVersion(cleanContent) + commentContent

def handlePodVersion(line):
    matches = re.findall(podNameFmt,line)
    if len(matches) > 0:
        name = matches[0][0]
    version = findVersion(name)
    if version == None:
        print("warning: dependency: %-60s find no version " % name)
        return line
    new_content = "pod '{library}', '={latest_version}'".format(library=name, latest_version=version)
    line = re.sub(podNameFmt, new_content, line)
    return line

def handlePodGit(line):
    matches = re.findall(podGitFmt,line)
    if len(matches) > 0:
        git = matches[0][1]
    tag = findTag(git)
    if tag == None:
        print("warning: git: %-60s find no tag " % git)
        return line
    match = re.search(r":tag\s*=>\s*'([^']+)'", line)  # 使用正则表达式匹配 ":tag => ''"
    if match:
        line = re.sub(r":tag\s*=>\s*'([^']+)'", ":tag => '%s'" % tag, line)  # 替换 ":tag =>"
    else:
        line += ", :tag => '%s'" % tag  # 在行末尾添加 ":tag =>"
    return line
    

def handleContent():
    replaceContent = []
    with open('Podfile', 'r') as file:
        fileContent = file.read()
    
    lines = fileContent.splitlines()
    for line in lines:
        replaceLine = handlLine(line)
        replaceContent.append(replaceLine)

    newContent = '\n'.join(replaceContent)
    with open('Podfile', 'w') as file:
        file.write(newContent) 

if __name__ == "__main__":
    handleContent()
