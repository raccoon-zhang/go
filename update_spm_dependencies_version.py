# 使用规范
# - 需要将脚本放到project.pbxproj同级目录执行
# - 只会替换固定版本号的库，对于其他类型的库不做处理

import re
import subprocess

def extractRepoInfo():
    #匹配依赖字符
    pattern = r'\/\* Begin XCRemoteSwiftPackageReference section \*\/([\s\S]+?)\/\* End XCRemoteSwiftPackageReference section \*\/'
    with open('project.pbxproj', 'r') as file:
        raw_string = file.read()
        matches = re.findall(pattern, raw_string, re.DOTALL)
    return matches[0]

def replaceVersion(replacements,originString):
    for package, version in replacements.items():
        if package in originString:
            pattern = r'(({})[\s\S]+?kind = exactVersion;[\s\S]+?version = )[\w.-]+;'.format(re.escape(package))
            replacement = r'\g<1>{};'.format(version)
            originString = re.sub(pattern, replacement, originString)
    return originString

def replaceSourceFileString(oldString, newString):
        with open('project.pbxproj', 'r') as file:
            fileString = file.read()
            newString = fileString.replace(oldString, newString)

        with open('project.pbxproj', 'w') as file:
            file.write(newString)

def findRepoName(line):
    pattern = r'.*?XCRemoteSwiftPackageReference "(.*?)"'
    matches = re.findall(pattern, line, re.DOTALL)
    if matches:
        return matches[0]
    else:
        return ""

def findRepoUrl(line):
    pattern = r'.*?repositoryURL = "(.*?)";'
    matches = re.findall(pattern, line, re.DOTALL)
    if matches:
        return matches[0]
    else:
        return ""

def findRepoKind(line):
    pattern = r'.*?kind = (.*?);'
    matches = re.findall(pattern, line, re.DOTALL)
    if matches:
        return matches[0]
    else:
        return ""


def getRepoAndNewVersion(originString):
    # 创建字典，存储库名称和version的对应关系
    libraries = {}
    lines = originString.splitlines()
    preToInsert = ["",""]
    for line in lines:
        if preToInsert[0] == "":
            preToInsert[0] = findRepoName(line)
        elif preToInsert[1] == "":
            preToInsert[1] = findRepoUrl(line)
        else:
            kind = findRepoKind(line)
            if kind == "":
                continue
            elif kind != "exactVersion":
                preToInsert = ["",""]
            else:
                latestVerison = getLatestVersion(preToInsert[1])
                if latestVerison:
                    libraries[preToInsert[0]] = latestVerison
                else:
                    print("warning: responsity:%-20s get version failed! url:%s" % (preToInsert[0],preToInsert[1]))
                preToInsert = ["",""]
    print()
    # 打印提取的结果
    for library_name, version in libraries.items():
        print(f"库名称: {library_name:<20} version: {version}")
    return libraries

def getLatestVersion(repoUrl):
    try:
        result = subprocess.check_output("git ls-remote --tags %s" % repoUrl, shell=True).decode("utf-8").split("\n")
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


if __name__=="__main__":
    repoInfo = extractRepoInfo()
    # 提取最新版本号
    versionDic = getRepoAndNewVersion(repoInfo)
    # 替换各个库的版本号
    newRepoInfo = replaceVersion(versionDic, repoInfo)
    # 替换源文件
    replaceSourceFileString(repoInfo, newRepoInfo)

