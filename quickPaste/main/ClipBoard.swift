//
//  ClipBoard.swift
//  quickPaste
//
//  Created by zhangshaohua on 2024/9/25.
//

import Foundation
import Cocoa

class ClipBoard: NSObject {
    static let shared = ClipBoard()
    struct HookInfo {
        let identifier: String
        let hook: Hook
    }
    typealias Hook = (String) -> String

    private let pasteboard = NSPasteboard.general
    private let timerInterval = 1.0

    private var changeCount: Int
    private var hooks: [HookInfo]
    public var enabled:Bool = true

    override init() {
        changeCount = pasteboard.changeCount
        hooks = []
        super.init()
        onNewCopy("default", underlineToHump(_:))
    }
    
    func onNewCopy(_ identifier:String, _ hook: @escaping Hook) {
        hooks.removeAll(where: { $0.identifier == identifier })
        hooks.append(HookInfo(identifier: identifier, hook: hook))
    }
    
    func onNewCopyFirst(_ identifier:String, _ hook: @escaping Hook) {
        hooks.removeAll(where: { $0.identifier == identifier })
        hooks.insert(HookInfo(identifier: identifier, hook: hook), at: 0)
    }

    func startListening() {
      Timer.scheduledTimer(timeInterval: timerInterval,
                           target: self,
                           selector: #selector(checkForChangesInPasteboard),
                           userInfo: nil,
                           repeats: true)
    }

    func copy(_ string: String) {
      pasteboard.declareTypes([NSPasteboard.PasteboardType.string], owner: nil)
      pasteboard.setString(string, forType: NSPasteboard.PasteboardType.string)
    }

    @objc func checkForChangesInPasteboard() {
      guard pasteboard.changeCount != changeCount, enabled else {
        return
      }

      if var lastItem = pasteboard.string(forType: NSPasteboard.PasteboardType.string) {
        for hook in hooks {
            lastItem = hook.hook(lastItem)
            copy(lastItem)
        }
      }

      changeCount = pasteboard.changeCount
    }
    
    /// 下划线字符串转换为驼峰
    func underlineToHump(_ string: String) -> String {
        guard string.isUnderline() else {
            return string
        }
        var result = ""
        let list = string.components(separatedBy: "_")
        for (index, item) in list.enumerated() {
            if index == 0 {
                result.append(item)
            } else {
                result.append(item.capitalized)
            }
        }
        return result
    }
}

extension String {
    /// 检查是否是下划线命名
    func isUnderline() -> Bool {
        var hasUnderline = false
        // 检查是否以字母或数字开头
        let firstCharacter = self[startIndex]
        guard (firstCharacter.isLetter || firstCharacter.isNumber), firstCharacter.isASCII else { return false }
        
        for item in self {
            /// 如果不是数字，下划线和字母则认为是非驼峰命名
            if !(item.isNumber ||  item.isLetter || [".", " ", "_"].contains(item)) || !item.isASCII{
                return false
            }
            if item == "_" {
                hasUnderline = true
            }
        }
        return hasUnderline
    }

}
