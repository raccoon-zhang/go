//
//  Label.swift
//  quickPaste
//
//  Created by zhangshaohua on 2024/9/26.
//

import Foundation
import Cocoa

class Label: NSTextField {
    
    init() {
        super.init(frame: .zero)
        isBezeled = false
        drawsBackground = false
        isEditable = false
        isSelectable = false
        font = NSFont.systemFont(ofSize: 12)
        textColor = NSColor.black
    }
    
    required init?(coder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }
}
