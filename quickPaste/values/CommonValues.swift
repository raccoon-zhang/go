//
//  CommonValues.swift
//  quickPaste
//
//  Created by zhangshaohua on 2024/9/30.
//

import Foundation
import Cocoa

class CommonValues: NSObject {
    static let shared = CommonValues()
    public var prefix = ""
    public var statusItem: NSStatusItem  {
        NSStatusBar.system.statusItem(withLength: NSStatusItem.variableLength)
    }
}
